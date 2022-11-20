import { UserAccount, Feed, DiscordChannel } from '../types/user'
import { FeedRow, NewFeedRow } from './feed-row'
import { useState, useEffect } from 'react'
import { useAuth0 } from "@auth0/auth0-react";
import { setConstantValue } from 'typescript';

export const UserModal = () => {
  const {
    getAccessTokenSilently,
  } = useAuth0();

  const [user, setUser] = useState<UserAccount>({ userID: -1, username: '', feedList: [], channelList: [] })

  //const [userId, setUserId] = useState(-1)
  //const [username, setUsername] = useState('')
  //const [feeds, setFeeds] = useState<Feed[]>([])
  //const [channels, setChannels] = useState<DiscordChannel[]>([])
  
  useEffect(() => {
    const url = process.env.REACT_APP_APIGW_ENDPOINT!
    const fetchUser = async (id: number) => {
      const accessToken = await getAccessTokenSilently()
      const resp = await fetch(`${url}user?userID=${id}`, {
        headers: {
          authorization: `Bearer ${accessToken}`
        },
      })
        .then(res => { return res.json() })
        .then(data => { return JSON.parse(data["body"]) })

      setUser(resp)
    }

    fetchUser(2)
  }, [])

  if (!user) {
    return (
      <div>
        <div>Loading your user account ...</div>
      </div>
    )
  }

  const accountRows = (feedList: Feed[]) => feedList.map((feed: Feed) => {
    return (
      <FeedRow key={`FeedRow-${feed.feedID}`} feed={feed} />
    )
  })

  const submitForm = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    const newFeed: Feed = {
      feedID: user.feedList.length + 1,
      title: (event.currentTarget.elements.namedItem('feedName') as HTMLInputElement).value,
      url: (event.currentTarget.elements.namedItem('feedUrl') as HTMLInputElement).value,
      timeFormat: '123'
    }

    const updatedUser: UserAccount = Object.assign({}, user)
    updatedUser.feedList = updatedUser.feedList.concat(newFeed)
    
    console.log(`Submitted form: \n name: ${newFeed.title}\n url: ${newFeed.url}\n id: ${newFeed.feedID}\n timeformat: ${newFeed.timeFormat}`)
    setUser(updatedUser)
  }

  return (
    <div className="overflow-hidden border rounded-lg">
        <form onSubmit={event => submitForm(event)}>
          <table className="divide-y divide-gray-200">
            <thead
              key={`FeedTableHeader-${user.userID}`} 
            >
              <tr>
                <td
                  className="w-1/12 py-3 text-xs font-bold text-left text-gray-500 uppercase"
                >
                  Feed ID
                </td>

                <td
                  className="w-4/12 px-6 py-3 text-xs font-bold text-left text-gray-500 uppercase"
                >
                  Feed Name
                </td>

                <td
                  className="w-3/12 px-6 py-3 text-xs font-bold text-left text-gray-500 uppercase"
                >
                  Feed URL
                </td>

                <td
                  className="w-3/12 px-6 py-3 text-xs font-bold text-left text-gray-500 uppercase"
                >
                  Feed Timestamp
                </td>
                <td
                  className="w-1/12 py-3 text-xs font-bold text-left text-gray-500 uppercase"
                >
                  &nbsp;
                </td>
              </tr>
            </thead>
            <tbody key={`FeedTableBody-${user.userID}`} className="divide-y divide-gray-200">
              {user?.feedList.map((feed: Feed) => {
                return <FeedRow key={`FeedRow-${feed.feedID}`} feed={feed} />
              })}
              <NewFeedRow />
            </tbody>
          </table>
        </form>
    </div>
  )
}