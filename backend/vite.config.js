import { defineConfig } from 'vite'
import { resolve } from 'path'
import { readdirSync } from 'fs'

/**
 * Recursively find all files with given extension in directory
 */
function findEntryFiles(dir, basePath = '', extension = '.css') {
  const entries = {}
  
  try {
    const items = readdirSync(dir, { withFileTypes: true })
    
    items.forEach(item => {
      const fullPath = resolve(dir, item.name)
      const relativePath = basePath ? `${basePath}/${item.name}` : item.name
      
      if (item.isDirectory()) {
        Object.assign(entries, findEntryFiles(fullPath, relativePath, extension))
      } else if (item.name.endsWith(extension)) {
        // Remove extension from entry name
        const name = relativePath.replace(extension, '')
        entries[name] = fullPath
      }
    })
  } catch (err) {
    // Directory doesn't exist yet, that's okay
  }
  
  return entries
}

// Auto-discover all CSS and JS entry files
const cssEntries = findEntryFiles(resolve(__dirname, 'assets/css'), '', '.css')
const jsEntries = findEntryFiles(resolve(__dirname, 'assets/js'), '', '.js')

// Combine entries
const allEntries = { ...cssEntries, ...jsEntries }

export default defineConfig({
  build: {
    outDir: 'static/dist',
    emptyOutDir: true,
    rollupOptions: {
      // Only set input if we have entries, otherwise use empty object
      input: Object.keys(allEntries).length > 0 ? allEntries : undefined,
      output: {
        // JS files go to their paths
        entryFileNames: '[name].js',
        // CSS files keep their paths
        assetFileNames: '[name].css'
      }
    },
    minify: true,
    cssMinify: true,
  }
})
