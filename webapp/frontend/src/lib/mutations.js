import { gql } from "@apollo/client";

export const LOGIN = gql`
  mutation Login($email: String!, $password: String!) {
    login(email: $email, password: $password)
  }
`;

export const SIGNUP = gql`
  mutation Signup($email: String!, $password: String!) {
    signup(email: $email, password: $password)
  }
`;

export const ME_QUERY = gql`
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

export const PODCAST_QUERY = gql`
  query GetPodcast($date: String) {
    podcast(date: $date) {
      date
      url
    }
  }
`;

export const UPDATE_PREFS = gql`
  mutation UpdatePreferences($country: String!, $topic: String!) {
    updatePreferences(country: $country, topic: $topic) {
      country
      topic
    }
  }
`;