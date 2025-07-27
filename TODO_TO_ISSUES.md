# Converting README TODOs to GitHub Issues

This document explains how to convert the TODO items from the README.md into proper GitHub issues.

## Overview

The repository contained 8 TODO items in the README.md that have been converted into structured issue templates:

1. **Secure client UUIDs** - Security enhancement for client identification
2. **WebSocket security** - Authentication and security for WebSocket connections  
3. **Sprite animations** - Basic sprite animation system
4. **Dynamic animations** - Event-driven animation system
5. **Server self-cleanup** - Automatic resource management and cleanup
6. **Multiplayer lobbies** - Game room and lobby system
7. **UI** - User interface and HUD implementation
8. **Player customization** - Avatar and player personalization features

## Generated Files

The `scripts/convert-todos-to-issues.sh` script has generated:
- `issues_to_create/` directory with individual issue templates
- Each issue includes: title, description, acceptance criteria, labels, priority, and effort estimates
- `issues_to_create/README.md` with detailed usage instructions

## How to Create the Issues

### Option 1: GitHub CLI (Automated)
If you have [GitHub CLI](https://cli.github.com/) installed:

```bash
cd issues_to_create
for file in *.md; do
  if [ "$file" != "README.md" ]; then
    title=$(grep '# ' $file | sed 's/# //')
    body=$(tail -n +3 $file)
    labels=$(grep 'Labels' $file | sed 's/.*`\(.*\)`.*/\1/' | tr ',' ' ')
    gh issue create --title "$title" --body "$body" --label "$labels"
  fi
done
```

### Option 2: Manual Creation via GitHub Web Interface
1. Navigate to the [Issues page](https://github.com/plyr4/go-ebiten-multiplayer/issues)
2. Click "New issue"
3. For each file in `issues_to_create/`:
   - Copy the title (line starting with `#`)
   - Copy the content below the title
   - Add the suggested labels
   - Create the issue

### Option 3: GitHub API (Advanced)
Use the GitHub REST API to create issues programmatically if you have appropriate permissions.

## Next Steps

After creating the issues:

1. **Update README.md**: Replace the TODO section with references to the created issues
2. **Add issue numbers**: Update any documentation that references these features
3. **Organize with milestones**: Group related issues into milestones if desired
4. **Assign priorities**: Use GitHub's priority labels to organize the backlog
5. **Clean up**: Remove the `issues_to_create/` directory after issues are created

## Issue Structure

Each generated issue includes:
- **Clear title** describing the feature/task
- **Detailed description** explaining the purpose and context
- **Acceptance criteria** with specific, actionable checkboxes
- **Suggested labels** for categorization
- **Priority level** (High/Medium/Low)
- **Effort estimate** (Small/Medium/Large)
- **Traceability** back to the original TODO item

## Labels Used

The conversion script suggests these labels:
- `security` - Security-related issues
- `enhancement` - Feature enhancements
- `websocket` - WebSocket-specific functionality
- `graphics` - Visual/graphics features
- `animation` - Animation-related features
- `server` - Server-side functionality
- `performance` - Performance improvements
- `maintenance` - Maintenance and cleanup tasks
- `multiplayer` - Multiplayer-specific features
- `feature` - New features
- `ui` - User interface elements
- `frontend` - Client-side functionality
- `player` - Player-related features
- `customization` - Customization features

You can adjust or add labels as needed when creating the issues.