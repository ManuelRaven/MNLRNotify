/**
 * This script copies Bootswatch theme files from node_modules to the public folder
 * so they are available in production builds.
 */
import fs from 'fs';
import path from 'path';

const __dirname = path.resolve(); // Get the current directory

// Define source and destination directories
const sourcePath = path.resolve(__dirname, 'node_modules/bootswatch/dist');
const destPath = path.resolve(__dirname, 'frontend/public/themes');

// Make sure the destination directory exists
if (!fs.existsSync(destPath)) {
  fs.mkdirSync(destPath, { recursive: true });
}

// Read all theme directories from bootswatch
const themes = fs.readdirSync(sourcePath);

// Copy each theme's bootstrap.min.css file
themes.forEach(theme => {
  const sourceFile = path.join(sourcePath, theme, 'bootstrap.min.css');
  const destDir = path.join(destPath, theme);
  const destFile = path.join(destDir, 'bootstrap.min.css');
  
  // Skip if not a directory or doesn't contain bootstrap.min.css
  if (!fs.existsSync(sourceFile)) {
    console.log(`Skipping ${theme} - bootstrap.min.css not found`);
    return;
  }
  
  // Create theme subdirectory in public/themes if it doesn't exist
  if (!fs.existsSync(destDir)) {
    fs.mkdirSync(destDir, { recursive: true });
  }
  
  // Copy the file
  fs.copyFileSync(sourceFile, destFile);
  console.log(`Copied ${theme} theme to public directory`);
});

console.log('All themes copied successfully');
