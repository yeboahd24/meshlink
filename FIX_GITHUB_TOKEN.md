# Fix GitHub Token Issue

## Problem:
Your Personal Access Token lacks `workflow` scope to create GitHub Actions.

## Solution 1: Update Token Scope (Recommended)

### Steps:
1. Go to GitHub → Settings → Developer settings → Personal access tokens
2. Find your current token
3. Click "Edit" or create new token
4. Check the `workflow` scope checkbox
5. Update token
6. Try push again

## Solution 2: Push Without Workflow First

### Remove workflow temporarily:
```bash
# Move workflow out of the way
mv .github/workflows/build.yml build.yml.backup
rm -rf .github

# Push code first
git add .
git commit -m "Add MeshLink code without workflow"
git push origin development

# Add workflow later via GitHub web interface
```

## Solution 3: Create Workflow via GitHub Web

### Steps:
1. Push code without `.github` folder
2. Go to GitHub repository → Actions tab
3. Click "New workflow" → "Set up a workflow yourself"
4. Copy content from `build.yml.backup`
5. Commit directly on GitHub

## Quick Fix:
```bash
# Remove workflow for now
rm -rf .github
git add .
git commit -m "Remove workflow temporarily"
git push origin development

# Add workflow later with proper token
```

Then update your token scope and add the workflow back.