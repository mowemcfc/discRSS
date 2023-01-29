import { UserAccount, Feed, DiscordChannel, NewFeedParams } from '../types/user'
import { FeedRow, NewFeedRow } from './feed-row'
import React, { useState, useEffect } from 'react'
import { useAuth0 } from "@auth0/auth0-react";
import { handleErrors } from '../utils';
import { arrayBuffer } from 'stream/consumers';


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
  }, [feedList])

  const addNewFeedHandler = async (newFeedTitle: string, newFeedUrl: string) => {
    const newFeedParams: NewFeedParams = {
        title: newFeedTitle,
        url: newFeedUrl,
    }

    const accessToken = await getAccessTokenSilently()
    const resp = await fetch(
      `${process.env.REACT_APP_APIGW_ENDPOINT!}user/${userId.toString()}/feed`, {
      method: 'POST',
      body: JSON.stringify(newFeedParams),
      headers: {
        authorization: `Bearer ${accessToken}`
      },
    })
      .then(handleErrors)
      .then(res => res.json())
      .then(data => data.Body)

    setFeedListState([...feedListState, (resp as Feed)])
  }

  const removeFeedHandler = async (feedId: number) => {
    setFeedListState(
      feedListState.filter(el => el.feedId != feedId)
    )
  }

  return (
    <div className="overflow-hidden grid place-items-center border rounded-lg">
          <table className="divide-y divide-gray-200">
            <thead>
              <tr>
                <td
                  className="py-3 text-xs font-bold text-left text-gray-500 uppercase"
                >
                  Feed ID
                </td>
                <td
                  className="px-6 py-3 text-xs font-bold text-left text-gray-500 uppercase"
                >
                  Feed Name
                </td>
                <td
                  className="w-24 px-6 py-3 text-xs font-bold text-left text-gray-500 uppercase"
                >
                  Feed URL
                </td>
                <td
                  className="w-3 px-6 py-3 text-xs font-bold text-left text-gray-500 uppercase"
                >
                  Feed Timestamp
                </td>
                <td
                  className="py-3 text-xs font-bold text-left text-gray-500 uppercase"
                >
                  &nbsp;
                </td>
                <td
                  className="py-3 text-xs font-bold text-left text-gray-500 uppercase"
                >
                  &nbsp;
                </td>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200">
                {feedListState.map((feed: Feed) => {
                  return <FeedRow removalHandler={removeFeedHandler} feed={feed} />
                })}
                <NewFeedRow addNewFeedStateHandler={addNewFeedHandler} />
            </tbody>
          </table>
    </div>
  )
}