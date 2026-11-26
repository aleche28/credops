"""Check the docs tree: ADR numbering, metadata, index, supersede links and relative links.

Usage:
    python scripts/check_docs.py           # check, exit 1 on problems
    python scripts/check_docs.py --write   # regenerate the ADR index, then check
"""

from __future__ import annotations

import argparse
import re
import sys
from dataclasses import dataclass
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
DOCS = ROOT / "docs"
ADR_DIR = DOCS / "adr"
ADR_INDEX = ADR_DIR / "README.md"
INDEX_START = "<!-- adr-index:start -->"
INDEX_END = "<!-- adr-index:end -->"

# Markdown files whose relative links are checked, besides everything under docs/.
EXTRA_LINKED_FILES = ["README.md", "CLAUDE.md"]
EXTRA_LINKED_GLOBS = [".claude/skills/**/*.md"]

STATUSES = {"Proposed", "Accepted", "Rejected", "Deprecated", "Superseded"}
ADR_FILE_RE = re.compile(r"^(\d{4})-[a-z0-9]+(?:-[a-z0-9]+)*\.md$")
ADR_REF_RE = re.compile(r"^ADR-(\d{4})$")
META_RE = re.compile(r"^- (Status|Date|Supersedes|Superseded by): (.+?)\s*$")
DATE_RE = re.compile(r"^\d{4}-\d{2}-\d{2}$")
LINK_RE = re.compile(r"(?<!!)\[[^\]]*\]\(([^)\s]+)(?:\s+\"[^\"]*\")?\)")
HEADING_RE = re.compile(r"^(#{1,6})\s+(.*?)\s*#*\s*$")
FENCE_RE = re.compile(r"^\s*(```|~~~)")


@dataclass
class Adr:
    number: int
    path: Path
    title: str
    meta: dict[str, str]


def strip_code(lines: list[str]) -> list[str]:
    """Blank out fenced code blocks and inline code so they aren't parsed as links."""
    out, in_fence = [], False
    for line in lines:
        if FENCE_RE.match(line):
            in_fence = not in_fence
            out.append("")
        elif in_fence:
            out.append("")
        else:
            out.append(re.sub(r"`[^`]*`", "", line))
    return out


def slugify(heading: str) -> str:
    """GitHub-style heading anchor."""
    text = re.sub(r"\[([^\]]*)\]\([^)]*\)", r"\1", heading)  # links -> link text
    text = re.sub(r"[`*_~]", "", text).strip().lower()
    text = re.sub(r"[^\w\- ]", "", text)
    return text.replace(" ", "-")


def anchors(path: Path) -> set[str]:
    seen: dict[str, int] = {}
    result = set()
    for line in strip_code(path.read_text(encoding="utf-8").splitlines()):
        m = HEADING_RE.match(line)
        if not m:
            continue
        slug = slugify(m.group(2))
        n = seen.get(slug, 0)
        seen[slug] = n + 1
        result.add(slug if n == 0 else f"{slug}-{n}")
    return result


def load_adrs(errors: list[str]) -> list[Adr]:
    adrs = []
    for path in sorted(ADR_DIR.glob("*.md")):
        if path.name in {"README.md", "template.md"}:
            continue
        rel = path.relative_to(ROOT)
        m = ADR_FILE_RE.match(path.name)
        if not m:
            errors.append(f"{rel}: file name must be NNNN-kebab-title.md")
            continue
        number = int(m.group(1))
        lines = path.read_text(encoding="utf-8").splitlines()
        title_m = re.match(rf"^# {m.group(1)}\. (.+)$", lines[0] if lines else "")
        if not title_m:
            errors.append(f"{rel}: first line must be '# {m.group(1)}. Title'")
            continue
        meta: dict[str, str] = {}
        for line in lines[1:]:
            if line.startswith("## "):
                break
            if mm := META_RE.match(line):
                meta[mm.group(1)] = mm.group(2)
        adrs.append(Adr(number, path, title_m.group(1), meta))
    return adrs


def check_adrs(adrs: list[Adr], errors: list[str]) -> None:
    numbers = [a.number for a in adrs]
    if numbers != list(range(1, len(adrs) + 1)):
        errors.append(f"docs/adr: ADR numbers must be 1..{len(adrs)} with no gaps, got {numbers}")
    by_number = {a.number: a for a in adrs}
    for adr in adrs:
        rel = adr.path.relative_to(ROOT)
        status = adr.meta.get("Status")
        if status not in STATUSES:
            errors.append(f"{rel}: Status must be one of {sorted(STATUSES)}, got {status!r}")
        if not DATE_RE.match(adr.meta.get("Date", "")):
            errors.append(f"{rel}: Date must be YYYY-MM-DD")
        if status == "Superseded" and "Superseded by" not in adr.meta:
            errors.append(f"{rel}: a Superseded ADR needs a 'Superseded by' line")
        if "Superseded by" in adr.meta and status != "Superseded":
            errors.append(f"{rel}: has 'Superseded by' but Status is {status}")
        for key, back in (("Supersedes", "Superseded by"), ("Superseded by", "Supersedes")):
            if key not in adr.meta:
                continue
            ref = ADR_REF_RE.match(adr.meta[key])
            other = by_number.get(int(ref.group(1))) if ref else None
            if other is None:
                errors.append(f"{rel}: '{key}' must name an existing ADR as ADR-NNNN")
            elif other.meta.get(back) != f"ADR-{adr.number:04d}":
                errors.append(
                    f"{rel}: '{key}: {adr.meta[key]}' but "
                    f"{other.path.relative_to(ROOT)} lacks '{back}: ADR-{adr.number:04d}'"
                )


def render_index(adrs: list[Adr]) -> str:
    rows = ["| ADR | Title | Status | Date |", "| --- | ----- | ------ | ---- |"]
    for a in adrs:
        rows.append(
            f"| [{a.number:04d}]({a.path.name}) | {a.title} "
            f"| {a.meta.get('Status', '')} | {a.meta.get('Date', '')} |"
        )
    return "\n".join([INDEX_START, *rows, INDEX_END])


def index_block(text: str) -> re.Match[str] | None:
    return re.search(re.escape(INDEX_START) + r".*?" + re.escape(INDEX_END), text, re.S)


def check_index(adrs: list[Adr], write: bool, errors: list[str]) -> None:
    rel = ADR_INDEX.relative_to(ROOT)
    text = ADR_INDEX.read_text(encoding="utf-8") if ADR_INDEX.exists() else ""
    block = index_block(text)
    if block is None:
        errors.append(f"{rel}: missing {INDEX_START} ... {INDEX_END} block")
        return
    expected = render_index(adrs)
    if block.group(0) == expected:
        return
    if write:
        ADR_INDEX.write_text(text[: block.start()] + expected + text[block.end() :], "utf-8")
        print(f"updated {rel}")
    else:
        errors.append(f"{rel}: ADR index is out of date; run scripts/check_docs.py --write")


def linked_files() -> list[Path]:
    files = set(DOCS.rglob("*.md"))
    files |= {ROOT / name for name in EXTRA_LINKED_FILES if (ROOT / name).exists()}
    for pattern in EXTRA_LINKED_GLOBS:
        files |= set(ROOT.glob(pattern))
    return sorted(files)


def check_links(errors: list[str]) -> None:
    anchor_cache: dict[Path, set[str]] = {}
    for path in linked_files():
        rel = path.relative_to(ROOT)
        lines = strip_code(path.read_text(encoding="utf-8").splitlines())
        for lineno, line in enumerate(lines, 1):
            for target in LINK_RE.findall(line):
                if re.match(r"^[a-z][a-z0-9+.-]*:", target):  # http:, mailto:, ...
                    continue
                file_part, _, anchor = target.partition("#")
                dest = (path.parent / file_part).resolve() if file_part else path
                if not dest.exists():
                    errors.append(f"{rel}:{lineno}: broken link {target}")
                    continue
                if anchor and dest.suffix == ".md":
                    if dest not in anchor_cache:
                        anchor_cache[dest] = anchors(dest)
                    if anchor not in anchor_cache[dest]:
                        errors.append(f"{rel}:{lineno}: missing anchor {target}")


def main(argv: list[str] | None = None) -> int:
    parser = argparse.ArgumentParser(description=__doc__.splitlines()[0])
    parser.add_argument("--write", action="store_true", help="regenerate the ADR index")
    args = parser.parse_args(argv)

    errors: list[str] = []
    adrs = load_adrs(errors)
    check_adrs(adrs, errors)
    check_index(adrs, args.write, errors)
    check_links(errors)

    for error in errors:
        print(error, file=sys.stderr)
    if errors:
        print(f"{len(errors)} problem(s) found", file=sys.stderr)
        return 1
    print(f"docs OK: {len(adrs)} ADRs, {len(linked_files())} files link-checked")
    return 0


if __name__ == "__main__":
    sys.exit(main())
