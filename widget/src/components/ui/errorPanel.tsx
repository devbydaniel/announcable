import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";

export function ErrorPanel() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Error</CardTitle>
      </CardHeader>
      <CardContent>
        <p>There was an error loading the release notes.</p>
      </CardContent>
    </Card>
  );
}
