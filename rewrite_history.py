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

def get_commit_info():
    """Get commit hashes and messages from oldest to newest"""
    result = subprocess.run(['git', 'log', '--format=%H|%s|%B', '--reverse'], 
                          capture_output=True, text=True)
    if result.returncode != 0:
        print(f"Error getting commit info: {result.stderr}")
        sys.exit(1)
    
    commits = []
    for line in result.stdout.strip().split('\n'):
        if '|' in line:
            parts = line.split('|', 2)
            if len(parts) >= 2:
                commits.append({
                    'hash': parts[0],
                    'message': parts[1],
                    'body': parts[2] if len(parts) > 2 else parts[1]
                })
    
    return commits

def rewrite_history_manually(commits, new_dates):
    """Rewrite history by creating new commits with new dates"""
    
    # First, create a temporary branch to work on
    subprocess.run(['git', 'checkout', '-b', 'temp-rewrite'], check=True)
    
    # Reset to the first commit's parent (empty state)
    first_commit = commits[0]['hash']
    result = subprocess.run(['git', 'rev-list', '--parents', first_commit], 
                          capture_output=True, text=True)
    
    # Start from empty repository state
    subprocess.run(['git', 'checkout', '--orphan', 'new-history'], check=True)
    subprocess.run(['git', 'rm', '-rf', '.'], check=False)  # Remove all files
    
    new_commit_hashes = []
    
    for i, commit_info in enumerate(commits):
        print(f"Rewriting commit {i+1}/{len(commits)}: {commit_info['hash'][:8]}")
        
        # Checkout the files from the original commit
        subprocess.run(['git', 'checkout', commit_info['hash'], '--', '.'], check=True)
        
        # Add all files
        subprocess.run(['git', 'add', '.'], check=True)
        
        # Set the environment variables for the new date
        git_date = format_date_for_git(new_dates[i])
        env = os.environ.copy()
        env['GIT_AUTHOR_DATE'] = git_date
        env['GIT_COMMITTER_DATE'] = git_date
        
        # Create the new commit
        parent_args = []
        if new_commit_hashes:  # If there's a previous commit, set it as parent
            parent_args = ['--parent', new_commit_hashes[-1]]
        
        try:
            result = subprocess.run([
                'git', 'commit', '-m', commit_info['message']
            ], env=env, capture_output=True, text=True, check=True)
            
            # Get the new commit hash
            result = subprocess.run(['git', 'rev-parse', 'HEAD'], 
                                  capture_output=True, text=True, check=True)
            new_commit_hash = result.stdout.strip()
            new_commit_hashes.append(new_commit_hash)
            
        except subprocess.CalledProcessError as e:
            print(f"Error creating commit: {e}")
            print(f"Stdout: {e.stdout}")
            print(f"Stderr: {e.stderr}")
            return False
    
    # Switch back to master and reset to the new history
    subprocess.run(['git', 'checkout', 'master'], check=True)
    subprocess.run(['git', 'reset', '--hard', new_commit_hashes[-1]], check=True)
    
    # Clean up temporary branches
    subprocess.run(['git', 'branch', '-D', 'temp-rewrite'], check=False)
    subprocess.run(['git', 'branch', '-D', 'new-history'], check=False)
    
    return True

def main():
    print("Git Commit History Rewriter (Manual Method)")
    print("===========================================")
    
    # Check for uncommitted changes
    result = subprocess.run(['git', 'status', '--porcelain'], capture_output=True, text=True)
    if result.stdout.strip():
        print("Warning: You have uncommitted changes. Stashing them temporarily...")
        subprocess.run(['git', 'stash', 'push', '-m', 'Temporary stash before history rewrite'], 
                      check=True)
        stashed = True
    else:
        stashed = False
    
    try:
        # Get commit information
        commits = get_commit_info()
        commit_count = len(commits)
        print(f"Found {commit_count} commits to rewrite")
        
        # Generate new dates
        print("Generating new dates between Aug 28, 2025 and Sep 14, 2025...")
        new_dates = generate_random_dates("2025-08-28", "2025-09-14", commit_count)
        
        # Show preview of date mapping
        print("\nDate mapping preview (first 5 commits):")
        for i in range(min(5, len(commits))):
            print(f"  {commits[i]['hash'][:8]} -> {format_date_for_git(new_dates[i])}")
        
        if len(commits) > 5:
            print(f"  ... and {len(commits) - 5} more commits")
        
        print(f"\nStarting history rewrite...")
        success = rewrite_history_manually(commits, new_dates)
        
        if success:
            print("✅ Successfully rewrote commit history!")
            return True
        else:
            print("❌ Failed to rewrite commit history")
            return False
            
    finally:
        # Restore stashed changes if any
        if stashed:
            print("Restoring stashed changes...")
            subprocess.run(['git', 'stash', 'pop'], check=True)

if __name__ == "__main__":
    if main():
        print("\n✅ Git commit dates successfully rewritten!")
        print("\nNew commit history (first 10 commits):")
        subprocess.run(['git', 'log', '--oneline', '--date=format:%Y-%m-%d %H:%M:%S', 
                       '--pretty=format:%h %ad %s', '-10'])
    else:
        print("\n❌ Failed to rewrite commit dates")
        sys.exit(1)
