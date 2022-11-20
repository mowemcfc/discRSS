import React from "react";
import PlusIcon from '../static/image/plus.png'
import { Feed } from '../types/user'


export const FeedRow = (props: {feed: Feed}) => {
  return (
      <tr key={props.feed.feedID}>
        <td key={props.feed.feedID} className="w-1/12 px-6 py-4 text-sm font-medium text-gray-800 whitespace-nowrap">{props.feed.feedID}</td>
        <td key={props.feed.title} className="w-3/12 px-6 py-4 text-sm font-medium text-gray-800 whitespace-nowrap">{props.feed.title}</td>
        <td key={props.feed.url} className="w-3/12 px-6 py-4 text-sm font-medium text-gray-800 whitespace-nowrap">{props.feed.url}</td>
        <td key={props.feed.timeFormat} className="w-3/12 px-6 py-4 text-sm font-medium text-gray-800 whitespace-nowrap">{props.feed.timeFormat}</td>
      </tr>
  );
};

export const NewFeedRow = () => {
  return (
      <tr>
          <td
            className="w-1/12 px-6 py-4"
          >
            &nbsp;
          </td>
          <td
            className="w-3/12 px-6 py-4"
          >
            <input type="text" placeholder="Feed Name"/>
          </td>
          <td
            className="w-3/12 px-6 py-4"
          >
            <input type="text" placeholder="Feed URL"/>
          </td>
          <td
            className="w-3/12 px-6 py-4"
          >
            &nbsp;
          </td>
          <td
            className="w-1/12 justify-start"
          >
            <input type="image" src={PlusIcon} width="20vw" height="20vw"/>
          </td>
      </tr>
  )
}
