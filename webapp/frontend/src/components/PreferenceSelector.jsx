import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Filter } from "lucide-react";

const COUNTRIES = [
  { value: 'ca', label: 'Canada' },
  { value: 'us', label: 'United States' },
  { value: 'uk', label: 'United Kingdom' },
  { value: 'au', label: 'Australia' },
  { value: 'de', label: 'Germany' },
  { value: 'fr', label: 'France' },
  { value: 'jp', label: 'Japan' },
];

const TOPICS = [
  { value: 'sports', label: 'Sports' },
  { value: 'technology', label: 'Technology' },
  { value: 'business', label: 'Business' },
  { value: 'science', label: 'Science' },
  { value: 'entertainment', label: 'Entertainment' },
  { value: 'health', label: 'Health & Wellness' },
  { value: 'politics', label: 'Politics' },
  { value: 'education', label: 'Education' },
];

const PreferenceSelector = ({ 
  country, 
  topic, 
  onCountryChange, 
  onTopicChange, 
  onApply,
  loading 
}) => {
  return (
    <Card className="bg-glass-bg border-glass-border backdrop-blur-sm">
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Filter className="h-5 w-5 text-primary" />
          Podcast Preferences
        </CardTitle>
        <CardDescription>
          Select your preferred country and topic to discover personalized AI-generated podcasts
        </CardDescription>
      </CardHeader>
      
      <CardContent className="space-y-4">
        <div className="grid md:grid-cols-2 gap-4">
          <div className="space-y-2">
            <label className="text-sm font-medium text-foreground">Country</label>
            <Select value={country} onValueChange={onCountryChange}>
              <SelectTrigger className="bg-background/50 border-border">
                <SelectValue placeholder="Select a country" />
              </SelectTrigger>
              <SelectContent className="bg-popover border-border">
                {COUNTRIES.map((c) => (
                  <SelectItem key={c.value} value={c.value}>
                    {c.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
          
          <div className="space-y-2">
            <label className="text-sm font-medium text-foreground">Topic</label>
            <Select value={topic} onValueChange={onTopicChange}>
              <SelectTrigger className="bg-background/50 border-border">
                <SelectValue placeholder="Select a topic" />
              </SelectTrigger>
              <SelectContent className="bg-popover border-border">
                {TOPICS.map((t) => (
                  <SelectItem key={t.value} value={t.value}>
                    {t.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        </div>
        
        <Button 
          onClick={onApply}
          variant="hero"
          className="w-full"
          disabled={!country || !topic || loading}
        >
          {loading ? 'Loading Podcasts...' : 'Find Podcasts'}
        </Button>
      </CardContent>
    </Card>
  );
};

export default PreferenceSelector;