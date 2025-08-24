import { gql } from "@apollo/client";
import { useQuery } from "@apollo/client/react";

const ME = gql`
  query Me {
    me {
      id
      email
      preferences {
        country
        topic
      }
    }
  }
`;

export default function Dashboard({ onLogout }) {
  const { data, loading, error } = useQuery(ME);

  if (loading) return <p>Loading...</p>;
  if (error) return <p>Error: {error.message}</p>;

  return (
    <div>
      <h2>Welcome, {data.me.email}</h2>
      <p>
        Podcast Preferences: {data.me.preferences.topic} from{" "}
        {data.me.preferences.country}
      </p>
      <button onClick={onLogout}>Logout</button>
    </div>
  );
}
