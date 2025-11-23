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

// Build a map to track which entries came from CSS vs JS sources
const sourceTypes = {}

// Combine entries: for files with both CSS and JS, give them unique entry names
const allEntries = {}

// Add CSS entries
Object.keys(cssEntries).forEach(key => {
  const entryName = `css__${key.replace(/\//g, '_')}`
  allEntries[entryName] = cssEntries[key]
  sourceTypes[entryName] = { type: 'css', path: key }
})

// Add JS entries
Object.keys(jsEntries).forEach(key => {
  const entryName = `js__${key.replace(/\//g, '_')}`
  allEntries[entryName] = jsEntries[key]
  sourceTypes[entryName] = { type: 'js', path: key }
})

export default defineConfig({
  build: {
    outDir: 'static/dist',
    emptyOutDir: true,
    rollupOptions: {
      // Only set input if we have entries, otherwise use empty object
      input: Object.keys(allEntries).length > 0 ? allEntries : undefined,
      output: {
        // JS files: use the original path from sourceTypes
        entryFileNames: (chunkInfo) => {
          const sourceInfo = sourceTypes[chunkInfo.name]
          if (sourceInfo && sourceInfo.type === 'js') {
            return `${sourceInfo.path}.js`
          }
          // For CSS sources or unknown entries, use default naming
          return '[name].js'
        },
        // CSS files: use the original path from sourceTypes
        assetFileNames: (assetInfo) => {
          if (assetInfo.name && assetInfo.name.endsWith('.css')) {
            // Find the matching source entry
            for (const [entryName, info] of Object.entries(sourceTypes)) {
              if (info.type === 'css' && assetInfo.name.includes(entryName)) {
                return `${info.path}.css`
              }
            }
          }
          return '[name][extname]'
        }
      }
    },
    minify: true,
    cssMinify: true,
  }
})
