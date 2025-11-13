import React from "react";
import Layout from "../components/layout";
import {
  Button,
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
  CardFooter,
  Badge,
  Alert,
  AlertTitle,
  AlertDescription,
  Separator,
  Input,
  Label,
  Textarea,
} from "../components/ui";

const ComponentShowcasePage: React.FC = () => {
  return (
    <Layout>
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="mb-12">
          <h1 className="text-4xl font-bold text-foreground mb-4">
            UI Component Showcase
          </h1>
          <p className="text-lg text-muted-foreground">
            A demonstration of our production-grade atom UI components built
            with Tailwind CSS and shadcn/ui principles.
          </p>
        </div>

        <Separator className="my-8" />

        {/* Buttons Section */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-6">Buttons</h2>
          <Card>
            <CardContent className="pt-6">
              <div className="space-y-4">
                <div>
                  <h3 className="text-sm font-medium mb-3">Variants</h3>
                  <div className="flex flex-wrap gap-3">
                    <Button variant="default">Default</Button>
                    <Button variant="secondary">Secondary</Button>
                    <Button variant="success">Success</Button>
                    <Button variant="destructive">Destructive</Button>
                    <Button variant="outline">Outline</Button>
                    <Button variant="ghost">Ghost</Button>
                    <Button variant="link">Link</Button>
                  </div>
                </div>
                <Separator />
                <div>
                  <h3 className="text-sm font-medium mb-3">Sizes</h3>
                  <div className="flex items-center flex-wrap gap-3">
                    <Button size="sm">Small</Button>
                    <Button size="default">Default</Button>
                    <Button size="lg">Large</Button>
                    <Button size="icon">🔍</Button>
                  </div>
                </div>
                <Separator />
                <div>
                  <h3 className="text-sm font-medium mb-3">States</h3>
                  <div className="flex flex-wrap gap-3">
                    <Button>Normal</Button>
                    <Button disabled>Disabled</Button>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        </section>

        {/* Badges Section */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-6">Badges</h2>
          <Card>
            <CardContent className="pt-6">
              <div className="flex flex-wrap gap-3">
                <Badge variant="default">Default</Badge>
                <Badge variant="secondary">Secondary</Badge>
                <Badge variant="success">Success</Badge>
                <Badge variant="destructive">Destructive</Badge>
                <Badge variant="warning">Warning</Badge>
                <Badge variant="outline">Outline</Badge>
              </div>
            </CardContent>
          </Card>
        </section>

        {/* Cards Section */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-6">Cards</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <Card>
              <CardHeader>
                <CardTitle>Simple Card</CardTitle>
                <CardDescription>
                  A card with header and content
                </CardDescription>
              </CardHeader>
              <CardContent>
                <p className="text-muted-foreground">
                  This is a simple card component that can be used to display
                  various types of content in a structured way.
                </p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Card with Footer</CardTitle>
                <CardDescription>
                  A complete card example
                </CardDescription>
              </CardHeader>
              <CardContent>
                <p className="text-muted-foreground">
                  This card includes a footer with action buttons.
                </p>
              </CardContent>
              <CardFooter className="flex gap-2">
                <Button variant="outline" className="flex-1">
                  Cancel
                </Button>
                <Button className="flex-1">Confirm</Button>
              </CardFooter>
            </Card>
          </div>
        </section>

        {/* Alerts Section */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-6">Alerts</h2>
          <div className="space-y-4">
            <Alert variant="default">
              <AlertTitle>Information</AlertTitle>
              <AlertDescription>
                This is a default alert for general information.
              </AlertDescription>
            </Alert>

            <Alert variant="success">
              <AlertTitle>🎉 Success!</AlertTitle>
              <AlertDescription>
                Your action was completed successfully.
              </AlertDescription>
            </Alert>

            <Alert variant="warning">
              <AlertTitle>⚠️ Warning</AlertTitle>
              <AlertDescription>
                Please review this information carefully before proceeding.
              </AlertDescription>
            </Alert>

            <Alert variant="destructive">
              <AlertTitle>❌ Error</AlertTitle>
              <AlertDescription>
                An error occurred. Please try again later.
              </AlertDescription>
            </Alert>
          </div>
        </section>

        {/* Form Components Section */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-6">
            Form Components
          </h2>
          <Card>
            <CardHeader>
              <CardTitle>Contact Form Example</CardTitle>
              <CardDescription>
                Form components with labels and inputs
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="name">Name</Label>
                <Input id="name" placeholder="Enter your name" />
              </div>
              <div className="space-y-2">
                <Label htmlFor="email">Email</Label>
                <Input
                  id="email"
                  type="email"
                  placeholder="Enter your email"
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="message">Message</Label>
                <Textarea
                  id="message"
                  placeholder="Type your message here"
                  rows={4}
                />
              </div>
            </CardContent>
            <CardFooter className="flex gap-2">
              <Button variant="outline">Clear</Button>
              <Button>Submit</Button>
            </CardFooter>
          </Card>
        </section>

        {/* Theme Colors Section */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-6">
            Theme Colors
          </h2>
          <Card>
            <CardContent className="pt-6">
              <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
                <div className="space-y-2">
                  <div className="h-20 bg-primary rounded-md flex items-center justify-center text-primary-foreground font-medium">
                    Primary
                  </div>
                </div>
                <div className="space-y-2">
                  <div className="h-20 bg-secondary rounded-md flex items-center justify-center text-secondary-foreground font-medium">
                    Secondary
                  </div>
                </div>
                <div className="space-y-2">
                  <div className="h-20 bg-success rounded-md flex items-center justify-center text-success-foreground font-medium">
                    Success
                  </div>
                </div>
                <div className="space-y-2">
                  <div className="h-20 bg-warning rounded-md flex items-center justify-center text-warning-foreground font-medium">
                    Warning
                  </div>
                </div>
                <div className="space-y-2">
                  <div className="h-20 bg-destructive rounded-md flex items-center justify-center text-destructive-foreground font-medium">
                    Destructive
                  </div>
                </div>
                <div className="space-y-2">
                  <div className="h-20 bg-muted rounded-md flex items-center justify-center text-muted-foreground font-medium">
                    Muted
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        </section>

        {/* Documentation Link */}
        <section>
          <Card className="bg-primary/5 border-primary/20">
            <CardContent className="pt-6">
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <h3 className="text-lg font-semibold text-foreground mb-2">
                    📚 Component Documentation
                  </h3>
                  <p className="text-muted-foreground mb-4">
                    Check out the complete documentation to learn how to use
                    these components in your project, including migration
                    guides and best practices.
                  </p>
                  <Button variant="default">
                    <a
                      href="https://github.com/steve-phan/egaldeutsch-go-lamda/blob/main/COMPONENT_LIBRARY.md"
                      target="_blank"
                      rel="noopener noreferrer"
                    >
                      View Documentation
                    </a>
                  </Button>
                </div>
              </div>
            </CardContent>
          </Card>
        </section>
      </div>
    </Layout>
  );
};

export default ComponentShowcasePage;

export const Head = () => <title>Component Showcase - EgalDeutsch</title>;
