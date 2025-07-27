# README TODOs Conversion

This directory contains the converted TODO items from the main README.md file, formatted as GitHub issues.

## Files
Each `.md` file represents one TODO item and contains:
- Issue title
- Detailed description  
- Acceptance criteria
- Suggested labels
- Priority level
- Effort estimate

## Usage

### Option 1: GitHub CLI (Automated)
If you have GitHub CLI installed:
```bash
cd issues_to_create
for file in *.md; do
  title=$(grep '# ' $file | sed 's/# //')
  body=$(tail -n +3 $file)
  labels=$(grep 'Labels' $file | sed 's/.*`\(.*\)`.*/\1/' | tr ',' ' ')
  gh issue create --title "$title" --body "$body" --label "$labels"
done
```

### Option 2: Manual Creation
1. Open each `.md` file
2. Copy the content
3. Create a new issue on GitHub
4. Paste the content and adjust as needed

## Next Steps
After creating the issues:
1. Update the main README.md to reference issue numbers instead of TODO items
2. Delete this directory (optional)
3. Close or update any related pull requests
