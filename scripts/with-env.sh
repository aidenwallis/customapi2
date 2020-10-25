  #!/bin/sh

base_dir="$(cd "$(dirname "$0")" >/dev/null 2>&1 && pwd)"
env_path="$base_dir/../.env"

if ! [ -f "$env_path" ]; then
	echo "[WARN] no .env file found. Executing without .env file"
	$"$@"
	exit
fi

env $(cat "$env_path" | xargs) "$@"
