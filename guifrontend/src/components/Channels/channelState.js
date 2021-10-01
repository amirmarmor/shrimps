import {Card, Form, ListGroup} from "react-bootstrap"
import React from "react"
import {useDispatch} from "react-redux"
import {actionAsync} from "../../features/config/configSlice"
import CheckBox from "./checkBox"

//TODO: get host from one place
const host = process.env["REACT_APP_BACKEND_HOST"] || "127.0.0.1"

function ChannelState(props) {
  const dispatch = useDispatch()
  const src = `http://${host}:8080/stream/${props.channel}`

  function handleChange(type, camera) {
    dispatch(actionAsync({type, channel: camera.toString()}))
  }

  function showIframe() {
    if (props.show) {
      return <iframe src={src} style={{width: "100%", height: "500px"}}/>
    }
  }

  let type = "checkbox"
  return (
    <Card>
      <Card.Header>
        <Card.Title as="h4">Channel {props.channel} State </Card.Title>
        <p className="card-category">24 Hours performance</p>
      </Card.Header>
      <Card.Body>
        <Form>
          <div key={`default-${type}`} className="mb-3">
            <ListGroup>
              <CheckBox
                checked={props.show}
                channel={props.channel}
                handleChange={handleChange}
                type={"show"}
              />
              <CheckBox
                checked={props.record}
                channel={props.channel}
                handleChange={handleChange}
                type={"record"}
              />
            </ListGroup>
          </div>
        </Form>
        {showIframe()}
      </Card.Body>
    </Card>
  )
}

export default ChannelState