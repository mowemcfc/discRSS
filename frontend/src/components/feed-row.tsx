import { ReactionForComment } from "aws-sdk/clients/codecommit";
import React, { useEffect, useState } from "react";
import PlusIcon from '../static/image/plus.png'
import RemoveIcon from '../static/image/remove.png'
import { Feed } from '../types/user'


interface FeedRowProps {
  feed: Feed
  removalHandler: (feedId: number) => void
}

export const FeedRow = ({feed, removalHandler}: FeedRowProps) => {

  const removeFeedHandler = (event: React.MouseEvent<HTMLInputElement>) => {
    removalHandler(feed.feedId)
  }

  return (
      <tr>
        <td  className="w-1/12 px-6 py-4 text-sm font-medium text-gray-800 whitespace-nowrap">{feed.feedId}</td>
        <td  className="w-3/12 px-6 py-4 text-sm font-medium text-gray-800 whitespace-nowrap">{feed.title}</td>
        <td  className="w-3/12 px-6 py-4 text-sm font-medium text-gray-800 whitespace-nowrap">{feed.url}</td>
        <td  className="w-2/12 px-6 py-4 text-sm font-medium text-gray-800 whitespace-nowrap">{feed.timeFormat}</td>
        <td  className="w-1/12 px-6 py-4 text-sm font-medium text-gray-800 whitespace-nowrap">
          <input type="image" onClick={removeFeedHandler} src={RemoveIcon} width="25vw" height="25vw"/>
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

  useEffect(() => {
    console.log(`new title: ${newFeedTitle}`)
    console.log(`new url: ${newFeedUrl}`)
    
  }, [newFeedTitle, newFeedUrl])

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
            <input type="text" onChange={newFeedTitleChangeHandler} 
              id="newFeedNameInput" name="feedName" placeholder="Feed Name" value={newFeedTitle}/>
          </td>
          <td
            className="px-6 py-4"
          >
            <input type="text" onChange={newFeedUrlChangeHandler} 
              id="newFeedUrlInput" name="feedUrl" placeholder="Feed URL" value={newFeedUrl}/>
          </td>
          <td
            className="px-6 py-4"
          >
            &nbsp;
          </td>
          <td
            className="w-1/12 px-6 py-4 justify-start"
          >
            <input type="image" onClick={newFeedSubmitHandler} src={PlusIcon} width="20vw" height="20vw"/>
          </td>
      </tr>
  )
}
