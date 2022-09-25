import { useAuth0 } from '@auth0/auth0-react';

export const useAuthnToken = async () => {
  const {
    getAccessTokenSilently
  } = useAuth0();

  const accessToken = await getAccessTokenSilently({
    audience:'https://mipjsyf7l2.execute-api.ap-southeast-2.amazonaws.com/v1/user',
    scope: 'read:user'
  })

  return accessToken;
}
