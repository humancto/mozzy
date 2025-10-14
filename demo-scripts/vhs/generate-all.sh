#!/bin/bash
# Generate all VHS demos at once

echo "🎬 Generating mozzy demo GIFs with VHS..."
echo ""

cd "$(dirname "$0")"

for tape in *.tape; do
  if [ -f "$tape" ]; then
    echo "📼 Recording: $tape"
    vhs "$tape"
    if [ $? -eq 0 ]; then
      echo "✅ Completed: $tape"
    else
      echo "❌ Failed: $tape"
    fi
    echo ""
  fi
done

echo ""
echo "🎉 All demos generated!"
echo "📁 Check the assets/ directory for GIFs"
ls -lh ../../assets/demo-*.gif 2>/dev/null || echo "⚠️  No GIFs found - check for errors above"
