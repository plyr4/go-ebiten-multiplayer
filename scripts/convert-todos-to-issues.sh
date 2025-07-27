#!/bin/bash

# Script to convert README TODOs to GitHub issues
# This script generates issue content that can be used to create GitHub issues

set -e

REPO_OWNER="plyr4"
REPO_NAME="go-ebiten-multiplayer"
OUTPUT_DIR="issues_to_create"

# Create output directory
mkdir -p "$OUTPUT_DIR"

echo "🎯 Converting README TODOs to GitHub Issues"
echo "Repository: $REPO_OWNER/$REPO_NAME"
echo "Output directory: $OUTPUT_DIR"
echo ""

# Define the TODO items with their details
declare -A todos

# TODO: secure client uuids
todos["secure-client-uuids"]="title:Secure client UUIDs
description:Currently client UUIDs may not be properly secured. We need to implement proper security measures to ensure client identification is secure and cannot be easily spoofed or manipulated.
labels:security,enhancement
acceptance_criteria:
- [ ] Review current UUID generation and handling
- [ ] Implement secure UUID validation on server side
- [ ] Add protection against UUID spoofing
- [ ] Document security measures implemented
priority:High
effort:Medium"

# TODO: websocket security
todos["websocket-security"]="title:Implement WebSocket security
description:WebSocket connections need proper security measures including authentication, authorization, and protection against common WebSocket vulnerabilities.
labels:security,websocket,enhancement
acceptance_criteria:
- [ ] Implement WebSocket authentication
- [ ] Add rate limiting for WebSocket connections
- [ ] Implement proper authorization checks
- [ ] Add protection against WebSocket-specific attacks
- [ ] Document security implementation
priority:High
effort:Large"

# TODO: sprite animations
todos["sprite-animations"]="title:Implement sprite animations
description:Add support for animated sprites to make the game more visually appealing and dynamic.
labels:graphics,animation,enhancement
acceptance_criteria:
- [ ] Design animation system architecture
- [ ] Implement sprite frame management
- [ ] Add animation timing controls
- [ ] Create example animated sprites
- [ ] Update rendering system to support animations
priority:Medium
effort:Medium"

# TODO: dynamic animations
todos["dynamic-animations"]="title:Add dynamic animations
description:Implement dynamic animations that can be triggered by game events and player actions, beyond static sprite animations.
labels:graphics,animation,enhancement
acceptance_criteria:
- [ ] Design dynamic animation system
- [ ] Implement event-driven animations
- [ ] Add smooth transitions and easing
- [ ] Create animation configuration system
- [ ] Test animations with multiplayer synchronization
priority:Medium
effort:Large"

# TODO: server self-cleanup
todos["server-self-cleanup"]="title:Implement server self-cleanup
description:The server should automatically clean up resources, remove inactive connections, and perform maintenance tasks to prevent resource leaks and ensure optimal performance.
labels:server,performance,maintenance
acceptance_criteria:
- [ ] Implement connection timeout and cleanup
- [ ] Add resource cleanup for disconnected clients
- [ ] Create periodic maintenance routines
- [ ] Add memory and resource monitoring
- [ ] Implement graceful shutdown procedures
priority:High
effort:Medium"

# TODO: multiplayer lobbies
todos["multiplayer-lobbies"]="title:Implement multiplayer lobbies
description:Add lobby system to allow players to create and join game rooms, making it easier to organize multiplayer sessions.
labels:multiplayer,feature,enhancement
acceptance_criteria:
- [ ] Design lobby system architecture
- [ ] Implement lobby creation and joining
- [ ] Add lobby listing and discovery
- [ ] Implement player management within lobbies
- [ ] Add lobby settings and configuration
- [ ] Create lobby UI components
priority:Medium
effort:Large"

# TODO: ui
todos["ui"]="title:Implement user interface
description:Create a proper user interface for the game including menus, HUD, settings, and other interactive elements.
labels:ui,frontend,enhancement
acceptance_criteria:
- [ ] Design UI/UX mockups
- [ ] Implement main menu system
- [ ] Add in-game HUD elements
- [ ] Create settings/options screen
- [ ] Implement responsive UI for different screen sizes
- [ ] Add accessibility features
priority:Medium
effort:Large"

# TODO: player customization
todos["player-customization"]="title:Add player customization
description:Allow players to customize their appearance, name, and other personal settings to enhance the multiplayer experience.
labels:player,customization,feature
acceptance_criteria:
- [ ] Design customization system
- [ ] Implement avatar/sprite customization
- [ ] Add player name and profile management
- [ ] Create customization UI
- [ ] Implement persistence of player settings
- [ ] Add multiplayer synchronization of customizations
priority:Low
effort:Medium"

# Function to create issue file
create_issue_file() {
    local issue_key="$1"
    local issue_data="$2"
    local filename="$OUTPUT_DIR/$issue_key.md"
    
    echo "Creating issue file: $filename"
    
    # Parse the issue data
    local title=$(echo "$issue_data" | grep "^title:" | cut -d: -f2- | sed 's/^ *//')
    local description=$(echo "$issue_data" | grep "^description:" | cut -d: -f2- | sed 's/^ *//')
    local labels=$(echo "$issue_data" | grep "^labels:" | cut -d: -f2- | sed 's/^ *//')
    local priority=$(echo "$issue_data" | grep "^priority:" | cut -d: -f2- | sed 's/^ *//')
    local effort=$(echo "$issue_data" | grep "^effort:" | cut -d: -f2- | sed 's/^ *//')
    
    # Extract acceptance criteria (everything between acceptance_criteria: and priority:)
    local acceptance_criteria=$(echo "$issue_data" | sed -n '/^acceptance_criteria:/,/^priority:/p' | sed '1d;$d')
    
    cat > "$filename" << EOF
# $title

## Description
$description

## Acceptance Criteria
$acceptance_criteria

## Labels
\`$labels\`

## Priority
$priority

## Estimated Effort
$effort

## Additional Context
This issue was created from a TODO item in the README.md file as part of converting the project roadmap into trackable GitHub issues.

---
**Original TODO:** \`- [ ] ${issue_key//-/ }\`
EOF
}

# Create issues for each TODO
echo "Creating issue files..."
for issue_key in "${!todos[@]}"; do
    create_issue_file "$issue_key" "${todos[$issue_key]}"
done

echo ""
echo "✅ Issue files created in $OUTPUT_DIR/"
echo ""
echo "📋 Next steps:"
echo "1. Review the generated issue files"
echo "2. Create GitHub issues using the provided content"
echo "3. Update README.md to reference the created issues"
echo ""
echo "💡 To create issues via GitHub CLI (if installed):"
echo "   cd $OUTPUT_DIR"
echo "   for file in *.md; do"
echo "     title=\$(grep '# ' \$file | sed 's/# //')"
echo "     body=\$(tail -n +3 \$file)"
echo "     labels=\$(grep 'Labels' \$file | sed 's/.*\`\\(.*\\)\`.*/\\1/' | tr ',' ' ')"
echo "     gh issue create --title \"\$title\" --body \"\$body\" --label \"\$labels\""
echo "   done"
echo ""
echo "📖 Manual creation:"
echo "   Use the content from each .md file to create issues through GitHub web interface"

# Create a summary file
cat > "$OUTPUT_DIR/README.md" << 'EOF'
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
EOF

echo "📄 Summary and instructions created in $OUTPUT_DIR/README.md"