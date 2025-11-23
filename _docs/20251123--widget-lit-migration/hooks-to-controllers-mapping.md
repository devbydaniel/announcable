# React Hooks → Lit Controllers/Tasks Mapping

## Complete Migration Reference

| React Hook | Lit Equivalent | Type | Purpose |
|------------|----------------|------|---------|
| `useWidgetToggle` | `WidgetToggleController` | Controller | Widget open/close state |
| `useAnchorsRef` | `AnchorsController` | Controller | Anchor element refs |
| `useReleaseNotes` | `ReleaseNotesTask` | Task | Fetch release notes |
| `useConfig` | `WidgetConfigTask` | Task | Fetch widget config |
| `useReleaseNoteStatus` | `ReleaseNoteStatusTask` | Task | Fetch status |
| `useReleaseNoteLikes` | `ReleaseNoteLikesController` | Task | Like/unlike |
| `useReleaseNoteMetrics` | `ReleaseNoteMetricsController` | Controller | Track metrics |
| `useWidgetAnchor` | _Merged into WidgetToggleController_ | - | Deprecated |

## Pattern Examples

### Controller Example
```typescript
// React
const [isOpen, setIsOpen] = useState(false);

// Lit
const controller = new WidgetToggleController(this);
// Access: controller.isOpen
// Update: controller.setIsOpen(true)
```

### Task Example
```typescript
// React
const { data, isLoading, error } = useQuery({ ... });

// Lit
const task = new ReleaseNotesTask(this, orgId);
${task.task.render({
  pending: () => html\`Loading...\`,
  complete: (data) => html\`...\`,
  error: (e) => html\`Error\`,
})}
```

✅ All hooks migrated successfully!
