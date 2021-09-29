import {ListGroupItem} from "react-bootstrap"
import React from "react"

function CheckBox(props) {
  let type = "checkbox"
  return (
    <ListGroupItem>
      <input
        name={`${props.type}-${props.channel}`}
        id={`${props.type}-${props.channel}`}
        type={type}
        checked={props.checked}
        onChange={() => props.handleChange(props.type, props.channel)}
      /> {props.type}
    </ListGroupItem>
  )
}

export default CheckBox