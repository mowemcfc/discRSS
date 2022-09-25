
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
    isLoading || !isAuthenticated ? null : 'https://qxxdo97bg8.execute-api.ap-southeast-2.amazonaws.com/v1/user',
    async (url) => {
      const accessToken = getAccessTokenSilently({
        scope: 'read:user'
      });
      const res = await fetch(`${url}?userID=${id}`, {
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