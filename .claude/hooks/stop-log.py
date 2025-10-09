#!/usr/bin/env python3
"""
Claude Code Stop hook script.
Saves a copy of the transcript when Claude finishes responding.

Usage: stop-log.py <project_root>
"""

import json
import sys
import shutil
from pathlib import Path
from datetime import datetime, timezone


def main():
    # Get project root from command line argument
    if len(sys.argv) < 2:
        print("✗ Error: Project root path required as first argument", file=sys.stderr)
        sys.exit(1)
    
    project_root = Path(sys.argv[1]).resolve()
    
    # Read the hook data from stdin
    hook_data = json.load(sys.stdin)
    
    # Extract relevant information
    session_id = hook_data.get("session_id")
    transcript_path = hook_data.get("transcript_path")
    stop_hook_active = hook_data.get("stop_hook_active", False)
    
    # Prevent infinite loops from stop hooks triggering more stop hooks
    if stop_hook_active:
        print("⚠ Stop hook already active, skipping", file=sys.stderr)
        sys.exit(0)
    
    # Generate ISO 8601 timestamp
    timestamp = datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
    
    # Create the log directory structure in the project root
    log_dir = project_root / ".ai-logs" / session_id
    log_dir.mkdir(parents=True, exist_ok=True)
    
    # Copy the transcript
    if transcript_path:
        transcript_source = Path(transcript_path).expanduser()
        if transcript_source.exists():
            transcript_dest = log_dir / f"{timestamp}-stop-transcript.log"
            shutil.copy2(transcript_source, transcript_dest)
            
            # Show relative path from project root for cleaner output
            try:
                rel_path = log_dir.relative_to(project_root)
                print(f"✓ Stop logged to {rel_path}/", file=sys.stderr)
            except ValueError:
                print(f"✓ Stop logged to {log_dir}/", file=sys.stderr)
        else:
            print(f"⚠ Transcript not found: {transcript_path}", file=sys.stderr)
    else:
        print("⚠ No transcript path provided", file=sys.stderr)


if __name__ == "__main__":
    try:
        main()
    except Exception as e:
        print(f"✗ Hook error: {e}", file=sys.stderr)
        sys.exit(1)
        