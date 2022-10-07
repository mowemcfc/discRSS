import React from "react";
import { Feed } from '../types/user'


export const FeedRow = (props: {feed: Feed}) => {
  return (
    <div> 
      <tr>
        <td>{props.feed.feedID}</td>
        <td>{props.feed.title}</td>
        <td>{props.feed.url}</td>
        <td>{props.feed.timeFormat}</td>
      </tr>
    </div>
  );
};
