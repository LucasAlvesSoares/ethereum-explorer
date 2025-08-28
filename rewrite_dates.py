#!/usr/bin/env python3

import subprocess
import random
from datetime import datetime, timedelta
import sys
import os

def generate_random_dates(start_date, end_date, count):
    """Generate random dates between start and end date, all between 4PM and 11:59PM"""
    dates = []
    start = datetime.strptime(start_date, "%Y-%m-%d")
    end = datetime.strptime(end_date, "%Y-%m-%d")
    
    # Generate dates across the range
    date_range = (end - start).days + 1
    
    for i in range(count):
        # Calculate which day this should be on (spread evenly but with some randomness)
        day_offset = int((i / count) * date_range) + random.randint(-1, 1)
        day_offset = max(0, min(day_offset, date_range - 1))
        
        target_date = start + timedelta(days=day_offset)
        
        # Random time between 4PM (16:00) and 11:59PM (23:59)
        hour = random.randint(16, 23)
        minute = random.randint(0, 59)
        second = random.randint(0, 59)
        
        full_datetime = target_date.replace(hour=hour, minute=minute, second=second)
        dates.append(full_datetime)
    
    # Sort dates to maintain chronological order
    dates.sort()
    return dates

def get_commit_hashes():
    """Get all commit hashes from oldest to newest"""
    result = subprocess.run(['git', 'log', '--format=%H', '--reverse'], 
                          capture_output=True, text=True)
    if result.returncode != 0:
        print(f"Error getting commit hashes: {result.stderr}")
        sys.exit(1)
    
    return [line.strip() for line in result.stdout.strip().split('\n') if line.strip()]

def format_date_for_git(dt):
    """Format datetime for Git (RFC 2822 format with timezone)"""
    # Format: "Wed Aug 28 18:41:07 2025 -0300"
    weekday = dt.strftime('%a')
    month = dt.strftime('%b')
    day = dt.day
    time = dt.strftime('%H:%M:%S')
    year = dt.year
    timezone = '-0300'  # Brazil timezone
    
    return f"{weekday} {month} {day} {time} {year} {timezone}"

def create_env_filter_script(commit_hashes, dates):
    """Create the environment filter script for git filter-branch"""
    script_content = "#!/bin/bash\n\ncase $GIT_COMMIT in\n"
    
    for i, commit_hash in enumerate(commit_hashes):
        if i < len(dates):
            git_date = format_date_for_git(dates[i])
            script_content += f'    {commit_hash})\n'
            script_content += f'        export GIT_AUTHOR_DATE="{git_date}"\n'
            script_content += f'        export GIT_COMMITTER_DATE="{git_date}"\n'
            script_content += f'        ;;\n'
    
    script_content += "    *)\n"
    script_content += "        # Keep original dates for any commits not in our list\n"
    script_content += "        ;;\n"
    script_content += "esac\n"
    
    return script_content

def main():
    print("Git Commit Date Rewriter")
    print("========================")
    
    # Get commit hashes
    commit_hashes = get_commit_hashes()
    commit_count = len(commit_hashes)
    print(f"Found {commit_count} commits to rewrite")
    
    # Generate new dates
    print("Generating new dates between Aug 28, 2025 and Sep 14, 2025...")
    new_dates = generate_random_dates("2025-08-28", "2025-09-14", commit_count)
    
    # Show preview of date mapping
    print("\nDate mapping preview (first 5 commits):")
    for i in range(min(5, len(commit_hashes))):
        print(f"  {commit_hashes[i][:8]} -> {format_date_for_git(new_dates[i])}")
    
    if len(commit_hashes) > 5:
        print(f"  ... and {len(commit_hashes) - 5} more commits")
    
    # Create environment filter script
    print("\nCreating environment filter script...")
    env_filter_content = create_env_filter_script(commit_hashes, new_dates)
    
    with open('date_filter.sh', 'w') as f:
        f.write(env_filter_content)
    
    os.chmod('date_filter.sh', 0o755)
    print("Created date_filter.sh")
    
    # Check for uncommitted changes
    result = subprocess.run(['git', 'status', '--porcelain'], capture_output=True, text=True)
    if result.stdout.strip():
        print("\nWarning: You have uncommitted changes. Stashing them temporarily...")
        subprocess.run(['git', 'stash', 'push', '-m', 'Temporary stash before date rewrite'], 
                      check=True)
        stashed = True
    else:
        stashed = False
    
    try:
        # Run git filter-branch
        print("\nRewriting commit dates...")
        env = os.environ.copy()
        env['FILTER_BRANCH_SQUELCH_WARNING'] = '1'
        
        result = subprocess.run([
            'git', 'filter-branch', '-f', '--env-filter', './date_filter.sh', 'HEAD'
        ], env=env, capture_output=True, text=True)
        
        if result.returncode != 0:
            print(f"Error running git filter-branch: {result.stderr}")
            print(f"Stdout: {result.stdout}")
            return False
        
        print("Successfully rewrote commit dates!")
        
        # Clean up git filter-branch refs
        try:
            subprocess.run(['git', 'for-each-ref', '--format=%(refname)', 'refs/original/'], 
                         capture_output=True, text=True, check=True)
            subprocess.run(['git', 'update-ref', '-d', 'refs/original/refs/heads/master'], 
                         check=True)
            print("Cleaned up filter-branch references")
        except subprocess.CalledProcessError:
            pass  # References might not exist
        
        return True
        
    finally:
        # Restore stashed changes if any
        if stashed:
            print("Restoring stashed changes...")
            subprocess.run(['git', 'stash', 'pop'], check=True)
        
        # Clean up temporary file
        if os.path.exists('date_filter.sh'):
            os.remove('date_filter.sh')
            print("Cleaned up temporary files")

if __name__ == "__main__":
    if main():
        print("\n✅ Git commit dates successfully rewritten!")
        print("\nNew commit history (first 10 commits):")
        subprocess.run(['git', 'log', '--oneline', '--date=format:%Y-%m-%d %H:%M:%S', 
                       '--pretty=format:%h %ad %s', '-10'])
    else:
        print("\n❌ Failed to rewrite commit dates")
        sys.exit(1)
