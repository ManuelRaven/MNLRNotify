import { defineStore } from 'pinia';
import { ref, watch } from 'vue';

export const availableThemes = [
  'brite',
  'cerulean',
  'cosmo',
  'cyborg',
  'darkly',
  'flatly',
  'journal',
  'litera',
  'lumen',
  'lux',
  'materia',
  'minty',
  'morph',
  'pulse',
  'quartz',
  'sandstone',
  'simplex',
  'sketchy',
  'slate',
  'solar',
  'spacelab',
  'superhero',
  'united', // This is the default as seen in main.ts
  'vapor',
  'yeti',
  'zephyr'
];

export const useThemeStore = defineStore('theme', () => {
  // Get theme from localStorage or use 'united' as default (matches import in main.ts)
  const theme = ref(localStorage.getItem('theme') || 'united');
  
  const changeTheme = async (newTheme: string) => {
    if (!availableThemes.includes(newTheme)) return;
    
    theme.value = newTheme;
    localStorage.setItem('theme', newTheme);
    
    // Dynamically load the selected theme
    try {
      await loadThemeCSS(newTheme);
    } catch (error) {
      console.error('Failed to load theme:', error);
    }
  };
    const loadThemeCSS = async (themeName: string) => {
    // Remove any previously loaded theme CSS
    const existingThemeLink = document.getElementById('bootswatch-theme');
    if (existingThemeLink) {
      existingThemeLink.remove();
    }
    
    // Create a new link element for the theme CSS
    const link = document.createElement('link');
    link.id = 'bootswatch-theme';
    link.rel = 'stylesheet';
    // Use public directory path which works in both dev and production
    link.href = `/themes/${themeName}/bootstrap.min.css`;
    document.head.appendChild(link);
    
    // Return a promise that resolves when the CSS is loaded
    return new Promise((resolve, reject) => {
      link.onload = () => resolve(true);
      link.onerror = () => reject(new Error(`Failed to load theme: ${themeName}`));
    });
  };
  
  // Initialize theme on store creation
  if (typeof window !== 'undefined') {
    loadThemeCSS(theme.value);
  }
  
  return {
    theme,
    availableThemes,
    changeTheme
  };
});
