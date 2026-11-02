def slugify(value: str) -> str:
    """Make a simple URL slug from words separated by whitespace."""
    return "-".join(value.lower().split())
