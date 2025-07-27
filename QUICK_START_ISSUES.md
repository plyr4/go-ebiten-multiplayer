# Quick Start: Creating GitHub Issues from TODO Items

## 🚀 Ready to Create Issues

All TODO items from README.md have been converted to structured issue templates in the `issues_to_create/` directory.

### One-Command Creation (GitHub CLI)

If you have [GitHub CLI](https://cli.github.com/) installed:

```bash
cd issues_to_create
for file in *.md; do
  if [ "$file" != "README.md" ]; then
    title=$(grep '# ' $file | sed 's/# //')
    body=$(tail -n +3 $file)
    labels=$(grep 'Labels' $file | sed 's/.*`\(.*\)`.*/\1/' | tr ',' ' ')
    echo "Creating issue: $title"
    gh issue create --title "$title" --body "$body" --label "$labels"
  fi
done
```

### Manual Creation Steps

1. Go to: https://github.com/plyr4/go-ebiten-multiplayer/issues/new
2. For each file in `issues_to_create/` (except README.md):
   - Copy the title (the `# Heading`)
   - Copy everything below the title
   - Add the suggested labels from the `Labels` section
   - Click "Submit new issue"

## 📋 8 Issues to Create

| Priority | Title | Labels | Effort |
|----------|-------|--------|--------|
| High | Secure client UUIDs | security, enhancement | Medium |
| High | Implement WebSocket security | security, websocket, enhancement | Large |
| High | Implement server self-cleanup | server, performance, maintenance | Medium |
| Medium | Implement sprite animations | graphics, animation, enhancement | Medium |
| Medium | Add dynamic animations | graphics, animation, enhancement | Large |
| Medium | Implement multiplayer lobbies | multiplayer, feature, enhancement | Large |
| Medium | Implement user interface | ui, frontend, enhancement | Large |
| Low | Add player customization | player, customization, feature | Medium |

## ✅ After Creating Issues

1. The issues will be automatically numbered (likely #2-#9)
2. You can optionally delete the `issues_to_create/` directory
3. Consider creating milestones to group related issues
4. Assign issues to team members as appropriate

## 📁 Files Created

- `scripts/convert-todos-to-issues.sh` - The conversion script
- `issues_to_create/*.md` - Individual issue templates  
- `TODO_TO_ISSUES.md` - Detailed documentation
- `.github/ISSUE_TEMPLATE/` - Templates for future issues
- Updated `README.md` - Now references issues instead of TODOs