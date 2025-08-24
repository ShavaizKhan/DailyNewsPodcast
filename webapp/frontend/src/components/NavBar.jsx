import { Button } from "@/components/ui/button";
import { Headphones, LogOut, User } from "lucide-react";
import { useNavigate } from "react-router-dom";
import { isAuthenticated, removeToken } from "@/lib/auth";

const NavBar = () => {
  const navigate = useNavigate();
  const authenticated = isAuthenticated();

  const handleLogout = () => {
    removeToken();
    navigate('/');
  };

  return (
    <nav className="fixed top-0 left-0 right-0 z-50 bg-glass-bg/80 backdrop-blur-md border-b border-glass-border">
      <div className="container mx-auto px-4 py-3 flex items-center justify-between">
        <div 
          className="flex items-center gap-2 cursor-pointer hover:scale-105 transition-transform"
          onClick={() => navigate('/')}
        >
          <Headphones className="h-8 w-8 text-primary" />
          <span className="text-xl font-bold text-foreground">PodcastAI</span>
        </div>
        
        <div className="flex items-center gap-4">
          {authenticated ? (
            <>
              <Button 
                variant="ghost" 
                onClick={() => navigate('/dashboard')}
                className="flex items-center gap-2"
              >
                <User className="h-4 w-4" />
                Dashboard
              </Button>
              <Button 
                variant="outline" 
                onClick={handleLogout}
                className="flex items-center gap-2"
              >
                <LogOut className="h-4 w-4" />
                Logout
              </Button>
            </>
          ) : (
            <>
              <Button 
                variant="ghost" 
                onClick={() => navigate('/login')}
              >
                Login
              </Button>
              <Button 
                variant="hero" 
                onClick={() => navigate('/signup')}
              >
                Sign Up
              </Button>
            </>
          )}
        </div>
      </div>
    </nav>
  );
};

export default NavBar;