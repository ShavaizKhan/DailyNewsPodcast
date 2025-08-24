import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Play, Clock, Globe } from "lucide-react";

const PodcastCard = ({ podcast, onPlay }) => {
  return (
    <Card className="group hover:scale-[1.02] transition-all duration-300 bg-glass-bg border-glass-border backdrop-blur-sm hover:shadow-glow cursor-pointer">
      <CardContent className="p-4">
        <div className="aspect-video bg-gradient-audio rounded-lg mb-4 flex items-center justify-center">
          {podcast.thumbnailUrl ? (
            <img 
              src={podcast.thumbnailUrl} 
              alt={podcast.title}
              className="w-full h-full object-cover rounded-lg"
            />
          ) : (
            <div className="text-4xl font-bold text-primary-foreground opacity-50">
              AI
            </div>
          )}
        </div>
        
        <h3 className="font-semibold text-foreground mb-2 line-clamp-2 group-hover:text-primary transition-colors">
          {podcast.title}
        </h3>
        
        <p className="text-sm text-muted-foreground mb-3 line-clamp-2">
          {podcast.description}
        </p>
        
        <div className="flex items-center justify-between text-xs text-muted-foreground mb-4">
          <div className="flex items-center gap-4">
            <div className="flex items-center gap-1">
              <Clock className="h-3 w-3" />
              {podcast.duration}
            </div>
            <div className="flex items-center gap-1">
              <Globe className="h-3 w-3" />
              {podcast.country.toUpperCase()}
            </div>
          </div>
          <span className="px-2 py-1 bg-primary/20 text-primary rounded text-xs font-medium">
            {podcast.topic}
          </span>
        </div>
        
        <Button 
          variant="audio" 
          onClick={() => onPlay(podcast)}
          className="w-full"
        >
          <Play className="h-4 w-4 mr-2" />
          Play Podcast
        </Button>
      </CardContent>
    </Card>
  );
};

export default PodcastCard;