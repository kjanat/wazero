"""Fetch selected spec tests from a pinned GitHub source archive."""

from pathlib import Path, PurePosixPath
import subprocess
import sys
import tarfile
import tempfile


def fetch(repo, ref, dest, sources):
    url = f"https://codeload.github.com/{repo}/tar.gz/{ref}"
    expected = {PurePosixPath(source): 0 for source in sources}
    fixtures = {}
    with tempfile.TemporaryDirectory() as tmp:
        archive = Path(tmp) / "source.tar.gz"
        subprocess.run(
            [
                "curl",
                "--fail",
                "--silent",
                "--show-error",
                "--location",
                "--retry",
                "3",
                "--connect-timeout",
                "15",
                "--max-time",
                "120",
                "--output",
                str(archive),
                url,
            ],
            check=True,
        )
        with tarfile.open(archive, "r:gz") as tar:
            for member in tar:
                path = PurePosixPath(*PurePosixPath(member.name).parts[1:])
                if not member.isfile() or path.suffix != ".wast":
                    continue
                for source in expected:
                    if path == source or path.parent == source:
                        expected[source] += 1
                        data = tar.extractfile(member).read()
                        if not data or path.name in fixtures:
                            raise ValueError(f"empty or duplicate fixture: {path}")
                        fixtures[path.name] = data
                        break

    missing = [str(source) for source, count in expected.items() if not count]
    if missing:
        raise ValueError(f"{repo}@{ref}: no fixtures found for {', '.join(missing)}")

    for name, data in fixtures.items():
        (dest / name).write_bytes(data)


if __name__ == "__main__":
    if len(sys.argv) < 5:
        sys.exit("usage: fetch-wast.py REPO REF DEST SOURCE [SOURCE ...]")
    repo, ref, dest, *sources = sys.argv[1:]
    try:
        fetch(repo, ref, Path(dest), sources)
    except (
        OSError,
        ValueError,
        tarfile.TarError,
        subprocess.CalledProcessError,
    ) as exc:
        sys.exit(f"fetch-wast: {exc}")
