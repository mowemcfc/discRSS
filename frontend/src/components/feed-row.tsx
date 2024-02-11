import React, { useEffect, useState } from "react";
import PlusIcon from '../static/image/plus.png'
import RemoveIcon from '../static/image/remove.png'
import { Feed } from '../types/user'


interface FeedRowProps {
  feedNumber: number
  feed: Feed
  removalHandler: (feedId: number) => void
}

export const FeedRow = ({feedNumber, feed, removalHandler}: FeedRowProps) => {

  const removeFeedHandler = (event: React.MouseEvent<HTMLInputElement>) => {
    removalHandler(feed.feedId)
  }

  return (
      <tr>
        <td  className="w-1/12 px-6 py-4">{feedNumber}</td>
        <td  className="w-3/12 px-6 py-4">{feed.title}</td>
        <td  className="w-3/12 px-6 py-4">{feed.url}</td>
        <td  className="w-1/12 px-6 py-4">
          <input type="image" alt="Remove feed" onClick={removeFeedHandler} src={RemoveIcon} width="25vw" height="25vw"/>
        </td>
      </tr>
  );
};

interface NewFeedRowProps {
  addNewFeedStateHandler: (newFeedTitle: string, newFeedUrl: string) => void
}

export const NewFeedRow = ({addNewFeedStateHandler}: NewFeedRowProps) => {
  const [ newFeedUrl, setNewFeedUrl ] = useState<string>('')
  const [ newFeedTitle, setNewFeedTitle ] = useState<string>('')

  useEffect(() => {}, [newFeedTitle, newFeedUrl])

  const newFeedTitleChangeHandler = (event:React.ChangeEvent<HTMLInputElement>) => {
    setNewFeedTitle(event.target.value)
  }
  const newFeedUrlChangeHandler = (event:React.ChangeEvent<HTMLInputElement>) => {
    setNewFeedUrl(event.target.value)
  }

  const newFeedSubmitHandler = (event: React.MouseEvent<HTMLInputElement>) => {
    addNewFeedStateHandler(newFeedTitle, newFeedUrl)
  }

  return (
      <tr>
          <td
            className="px-6 py-4"
          >
            &nbsp;
          </td>
          <td
            className="px-6 py-4"
          >
            <input type="text" className="indent-2 text-black/75" onChange={newFeedTitleChangeHandler} 
              id="newFeedNameInput" name="feedName" placeholder="..." value={newFeedTitle}/>
          </td>
          <td
            className="px-6 py-4"
          >
            <input type="text" className="indent-2 text-black/75" onChange={newFeedUrlChangeHandler} 
              id="newFeedUrlInput" name="feedUrl" placeholder="..." value={newFeedUrl}/>
          </td>
          <td
            className="w-1/12 px-6 py-4"
          >
            <input type="image" alt="Add feed" onClick={newFeedSubmitHandler} src={PlusIcon} width="20vw" height="20vw"/>
          </td>
      </tr>
  )
}
