import { UserAccount, Feed, DiscordChannel } from '../types/user'
import { FeedRow, NewFeedRow } from './feed-row'
import React, { useState, useEffect } from 'react'
import { useAuth0 } from "@auth0/auth0-react";
import { handleErrors } from '../utils';


interface FeedListProps {
  userId: number
  feedList: Feed[]
}

export const FeedList: React.FC<FeedListProps> = ({ userId, feedList }): JSX.Element => {
  const {
    getAccessTokenSilently,
  } = useAuth0();

  const [ feedListState, setFeedListState ] = useState<Feed[]>([])

  useEffect(() => {
    setFeedListState(feedList)
  })

  const submitForm = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    const newFeed: Feed = {
        feedId: feedList.length,
        title: (event.currentTarget.elements.namedItem('feedName') as HTMLInputElement).value,
        url: (event.currentTarget.elements.namedItem('feedUrl') as HTMLInputElement).value,
        timeFormat: '123' // TODO: remove
    }
    const newFeedParams = {
      userId: userId.toString(),
      newFeed: [ newFeed ]
    }

    console.log(`Submitted form: \n name: ${newFeed.title}\n url: ${newFeed.url}\n id: ${newFeed.feedId}\n timeformat: ${newFeed.timeFormat}`)

    const accessToken = await getAccessTokenSilently()
    const resp = await fetch(
      `${process.env.REACT_APP_APIGW_ENDPOINT!}user/feeds`, {
      method: 'POST',
      body: JSON.stringify(newFeedParams),
      headers: {
        authorization: `Bearer ${accessToken}`
      },
    })
      .then(handleErrors)
      .then(res => res.json())
      .then(data => JSON.parse(data?.["body"]))

    // There seems to be no simpler way to check object equality than using JSON.stringify()
    //  NOTE: this approach requires keys to be ordered the same way, as comparison is done against strings.
    console.log(JSON.stringify(resp))
    console.log(JSON.stringify(newFeedParams.newFeed))
    if (JSON.stringify(resp) === JSON.stringify(newFeedParams.newFeed)) {
      setFeedListState(feedListState.concat(newFeedParams.newFeed))
    } else {
      console.error('response did not match request body, indicating something went wrong')
    }
  }

  return (
    <div className="overflow-hidden border rounded-lg">
        <form onSubmit={event => submitForm(event)}>
          <table className="divide-y divide-gray-200">
            <thead>
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
            <tbody className="divide-y divide-gray-200">
              {feedListState.map((feed: Feed) => {
                return <FeedRow key={`FeedRow-${feed.feedId}`} feed={feed} />
              })}
              <NewFeedRow />
            </tbody>
          </table>
        </form>
    </div>
  )
}