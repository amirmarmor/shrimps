import {Col,Row} from "react-bootstrap"
import React from "react"
import ChannelState from "./channelState"

function ChannelRow(props) {
  return (
    <Row>
      <Col md={12}>
        <ChannelState
          show={props.show}
          record={props.record}
          channel={props.channel}
          key={`cahnnel-${props.channel}`}
        />
      </Col>
    </Row>
  )
}

export default ChannelRow