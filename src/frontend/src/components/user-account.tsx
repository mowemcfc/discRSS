
import { useAuthnToken } from "../hooks/useAuthnToken"
import { useAuth0 } from "@auth0/auth0-react"
import useSWR from "swr";

export const UserAccount = ({ id }: any) => {
  const {
    isAuthenticated,
    user,
    isLoading,
    getAccessTokenSilently
  } = useAuth0();

  const { data, error } = useSWR(
    isLoading || !isAuthenticated ? null : 'https://cbiobsxi12.execute-api.ap-southeast-2.amazonaws.com/prod/',
    async (url) => {
      const accessToken = await getAccessTokenSilently();
      const res = await fetch(`${url}user?userID=${id}`, {
        headers: {
          authorization: `Bearer ${accessToken}`
        },
      });
      return res.json();
    }
  );

  if (error) {
    return (
      <div> 
        There was an error loading your user account: {error.message}
      </div>
    )
  }

  if (!data) {
    return (
      <div>
        <h1>User account for {user?.name}</h1>
        <div>Loading your user account ...</div>
      </div>
    )
  }
  console.log(data)

  return (
      <div>
        <h1>User account for {user?.name}</h1>
        <div> {data.length} </div>
      </div>
  )
}