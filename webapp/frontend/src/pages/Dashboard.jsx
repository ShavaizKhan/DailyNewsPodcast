import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import NavBar from "@/components/NavBar";
import { useQuery, useLazyQuery, useMutation } from "@apollo/client/react";
import PreferenceSelector from "@/components/PreferenceSelector";
import PodcastCard from "@/components/PodcastCard";
import { isAuthenticated } from "@/lib/auth";
import { useToast } from "@/hooks/use-toast";
import { ME_QUERY, PODCAST_QUERY, UPDATE_PREFS } from "@/lib/mutations";

const Dashboard = () => {
  const navigate = useNavigate();
  const { toast } = useToast();

  const [country, setCountry] = useState("us");
  const [topic, setTopic] = useState("sports");
  const [currentPodcast, setCurrentPodcast] = useState(null);

  // Check authentication
  useEffect(() => {
    if (!isAuthenticated()) {
      navigate("/login");
    }
  }, [navigate]);

  // Load user
  const { data: meData, loading: meLoading } = useQuery(ME_QUERY, {
    onCompleted: (data) => {
      if (data?.me?.preferences) {
        setCountry(data.me.preferences.country);
        setTopic(data.me.preferences.topic);
      }
    },
    onError: () => {
      toast({
        title: "Error",
        description: "Failed to load user profile",
        variant: "destructive",
      });
    },
  });

 // Lazy query for podcast
  const [fetchPodcast, { data: podcastData, loading: podcastLoading }] =
    useLazyQuery(PODCAST_QUERY);

  // Mutation for preferences
  const [updatePreferences] = useMutation(UPDATE_PREFS, {
    onCompleted: (data) => {
      toast({
        title: "Preferences Saved",
        description: `Now showing podcasts about ${data.updatePreferences.topic} in ${data.updatePreferences.country.toUpperCase()}`,
      });
      fetchPodcast({ variables: { date: new Date().toISOString().split("T")[0] } });
    },
    onError: () => {
      toast({
        title: "Error",
        description: "Failed to save preferences",
        variant: "destructive",
      });
    },
  });

  const handleApplyPreferences = () => {
    updatePreferences({ variables: { country, topic } });
  };

  const handlePlayPodcast = () => {
    if (podcastData?.podcast) {
      setCurrentPodcast({
        title: `AI Podcast for ${topic} in ${country}`,
        duration: "Unknown",
        audioUrl: podcastData.podcast.url,
      });
      toast({
        title: "Now Playing",
        description: `Podcast for ${topic} in ${country}`,
      });
    } else {
      toast({
        title: "No podcast available",
        description: "Try again later or adjust preferences.",
      });
    }
  };

  return (
    <div className="min-h-screen bg-gradient-hero">
      <NavBar />
      
      <main className="pt-20 pb-8">
        <div className="container mx-auto px-4">
          <div className="mb-8">
            <h1 className="text-3xl font-bold text-foreground mb-2">
              Your AI Podcast Dashboard
            </h1>
            <p className="text-muted-foreground">
              Discover personalized AI-generated podcasts based on your preferences
            </p>
          </div>
          
          <div className="mb-8">
            <PreferenceSelector
              country={country}
              topic={topic}
              onCountryChange={setCountry}
              onTopicChange={setTopic}
              onApply={handleApplyPreferences}
              loading={podcastLoading}
            />
          </div>
          
          {podcastData?.podcast && (
            <div className="mb-8">
              <h2 className="text-2xl font-semibold text-foreground mb-4">
                Today’s Podcast
              </h2>
              <PodcastCard
                podcast={{
                  id: "today",
                  title: `AI Podcast for ${topic}`,
                  description: "Generated just for you.",
                  duration: "Unknown",
                  country,
                  topic,
                  audioUrl: podcastData.podcast.url,
                  createdAt: podcastData.podcast.date,
                }}
                onPlay={handlePlayPodcast}
              />
            </div>
          )}
          
          {currentPodcast && (
            <div className="fixed bottom-0 left-0 right-0 bg-glass-bg/95 backdrop-blur-md border-t border-glass-border p-4">
              <div className="container mx-auto flex items-center justify-between">
                <div className="flex items-center gap-4">
                  <div className="w-12 h-12 bg-gradient-audio rounded-lg flex items-center justify-center">
                    <span className="text-primary-foreground font-bold">AI</span>
                  </div>
                  <div>
                    <h4 className="font-medium text-foreground">{currentPodcast.title}</h4>
                    <p className="text-sm text-muted-foreground">{currentPodcast.duration}</p>
                  </div>
                </div>
                <audio controls src={currentPodcast.audioUrl} className="w-1/3" />
              </div>
            </div>
          )}
        </div>
      </main>
    </div>
  );
};

export default Dashboard;