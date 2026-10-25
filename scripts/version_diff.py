#!/usr/bin/env python3

"""Report native Terraform state schema version changes between two provider schemas.

Upjet rebuilds terraform.tfstate on every reconcile and stamps it with the
schema version baked in at code generation time, so Terraform never sees a
state older than the current provider and never runs a resource's state
upgrade function. A resource whose schema version moved therefore needs its
stored state migrated by hand; this script fails so that a bump cannot land
without someone deciding how.

usage: version_diff.py <generated resource list> <base JSON schema path> <bumped JSON schema path>
example usage: version_diff.py config/generated.lst .work/schema.json.1.0.1 config/schema.json
"""

import json
import sys


def resource_schemas(path):
    """Return the resource schemas of the single provider described by path."""
    with open(path) as f:
        schema = json.load(f)
    providers = schema.get("provider_schemas", {})
    if not providers:
        sys.exit(f"Cannot extract the provider name from the schema: {path}")
    # there is only ever one provider in a generated schema
    provider_name = next(iter(providers))
    return providers[provider_name].get("resource_schemas", {})


def main(list_path, base_path, bumped_path):
    print(f'Reporting schema changes between "{base_path}" as base version and "{bumped_path}" as bumped version')

    with open(list_path) as f:
        resources = json.load(f)
    base = resource_schemas(base_path)
    bumped = resource_schemas(bumped_path)

    changed = []
    missing = []
    for name in sorted(resources):
        if name not in base:
            # a newly generated resource has no stored state to migrate
            continue
        if name not in bumped:
            missing.append(name)
            continue
        base_version = base[name].get("version", 0)
        bumped_version = bumped[name].get("version", 0)
        if base_version != bumped_version:
            changed.append((name, base_version, bumped_version))

    # keep the report readable when stdout is a pipe and stderr is not
    sys.stdout.flush()
    for name in missing:
        print(f"{name} is generated but missing from the bumped schema", file=sys.stderr)
    sys.stderr.flush()

    if not changed:
        print(f"No native state schema version changes across {len(resources)} generated resources")
        return 0

    for name, base_version, bumped_version in changed:
        print(f"{name}:{base_version}-{bumped_version}")
    sys.stdout.flush()
    print(
        f"\n{len(changed)} resource(s) changed their native state schema version. Existing managed "
        "resources hold state written under the old version, and upjet stamps the new version into "
        "the state it rebuilds, so Terraform will not run the upstream state upgrade function. "
        "Migrate the stored state before merging this bump.",
        file=sys.stderr,
    )
    return 1


if __name__ == "__main__":
    if len(sys.argv) != 4:
        sys.exit(f"usage: {sys.argv[0]} <generated resource list> <base JSON schema path> <bumped JSON schema path>")
    sys.exit(main(sys.argv[1], sys.argv[2], sys.argv[3]))
