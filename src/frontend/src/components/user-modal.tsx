import { UserAccount, Feed, DiscordChannel } from '../types/user'
import { FeedRow, NewFeedRow } from './feed-row'
import React, { useState, useEffect } from 'react'
import { useAuth0 } from "@auth0/auth0-react";
import { setConstantValue } from 'typescript';

export const UserModal = () => {
  const {
    getAccessTokenSilently,
  } = useAuth0();

  const [user, setUser] = useState<UserAccount>({ userId: -1, username: '', feedList: [], channelList: [] })
  
  useEffect(() => {
    const fetchUser = async (id: number) => {
      const accessToken = await getAccessTokenSilently()
      const resp = await fetch(`${process.env.REACT_APP_APIGW_ENDPOINT!}user?userId=${id}`, {
        headers: {
          authorization: `Bearer ${accessToken}`
        },
      })
        .then(res => res.json())
        .then(data => JSON.parse(data["body"]))

      setUser(resp)
    }

    fetchUser(1) // TODO: get this value dynamically from auth0 ID 
  }, [])

  if (!user) {
    return (
      <div>
        <div>Loading your user account ...</div>
      </div>
    )
  }

  const submitForm = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    const newFeed: Feed = {
      feedId: user.feedList.length,
      title: (event.currentTarget.elements.namedItem('feedName') as HTMLInputElement).value,
      url: (event.currentTarget.elements.namedItem('feedUrl') as HTMLInputElement).value,
      timeFormat: '123' // TODO: remove
    }

    const updatedUser: UserAccount = Object.assign({}, user)
    updatedUser.feedList = updatedUser.feedList.concat(newFeed)
    
    console.log(`Submitted form: \n name: ${newFeed.title}\n url: ${newFeed.url}\n id: ${newFeed.feedId}\n timeformat: ${newFeed.timeFormat}`)

    const accessToken = await getAccessTokenSilently()
    const resp = await fetch(
      `${process.env.REACT_APP_APIGW_ENDPOINT!}user`, {
      method: 'POST',
      body: JSON.stringify(updatedUser),
      headers: {
        authorization: `Bearer ${accessToken}`
      },
    })
      .then(res => res.json())
      .then(data => JSON.parse(data?.["body"]))

    // There seems to be no simpler way to check object equality than using JSON.stringify()
    //  NOTE: this approach requires keys to be ordered the same way, as comparison is done against strings.
    if (JSON.stringify(resp) === JSON.stringify(updatedUser)) {
      setUser(updatedUser)
    } else {
      console.error('response did not match request body, indicating something went wrong')
    }
  }

  return (
    <div className="overflow-hidden border rounded-lg">
        <form onSubmit={event => submitForm(event)}>
          <table className="divide-y divide-gray-200">
            <thead
              key={`FeedTableHeader-${user.userId}`} 
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
            <tbody key={`FeedTableBody-${user.userId}`} className="divide-y divide-gray-200">
              {user.feedList.map((feed: Feed) => {
                return <FeedRow key={`FeedRow-${feed.feedId}`} feed={feed} />
              })}
              <NewFeedRow />
            </tbody>
          </table>
        </form>
    </div>
  )
}