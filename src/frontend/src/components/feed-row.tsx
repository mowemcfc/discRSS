import React from "react";
import { Feed } from '../types/user'


export const FeedRow = (props: {feed: Feed}) => {
  return (
      <tr key={props.feed.feedID}>
        <td key={props.feed.feedID} className="px-6 py-4 text-sm font-medium text-gray-800 whitespace-nowrap">{props.feed.feedID}</td>
        <td key={props.feed.title} className="px-6 py-4 text-sm font-medium text-gray-800 whitespace-nowrap">{props.feed.title}</td>
        <td key={props.feed.url} className="px-6 py-4 text-sm font-medium text-gray-800 whitespace-nowrap">{props.feed.url}</td>
        <td key={props.feed.timeFormat} className="px-6 py-4 text-sm font-medium text-gray-800 whitespace-nowrap">{props.feed.timeFormat}</td>
      </tr>
  );
};
