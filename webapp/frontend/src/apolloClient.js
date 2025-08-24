import { ApolloClient, InMemoryCache, ApolloLink, HttpLink } from "@apollo/client";

// Http link for GraphQL endpoint
const httpLink = new HttpLink({
  uri: "http://localhost:8080/query", // Replace with your Go GraphQL server
});

// Auth link using ApolloLink
const authLink = new ApolloLink((operation, forward) => {
  const token = localStorage.getItem("token");
  if (token) {
    operation.setContext({
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
  }
  return forward(operation);
});

// Apollo Client
export const client = new ApolloClient({
  link: authLink.concat(httpLink),
  cache: new InMemoryCache(),
});
