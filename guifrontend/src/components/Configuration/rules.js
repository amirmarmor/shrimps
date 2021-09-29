import {Button, Col, Form, Row} from "react-bootstrap"
import React from "react"

function RuleRow(props) {
  return (
    <Row className={"align-items-bottom"} key={`rule-${props.rule.id}`}>
      <Col className={"pr-1"} md={2}>
        <Form.Group>
          <label>Type</label>
            <Form.Check
              label={"Video"}
              value={"video"}
              type={"radio"}
              id={`1-type-${props.id}`}
              checked={props.rule.type === "video"}
              name={"type"}
              onChange={e => props.handleRuleChange(e, props.id)}
            />
            <Form.Check
              label={"Image"}
              value={"image"}
              type={"radio"}
              id={`2-type-${props.id}`}
              checked={props.rule.type === "image"}
              name={"type"}
              onChange={e => props.handleRuleChange(e, props.id)}
            />
        </Form.Group>
      </Col>

      <Col className="ruleIdpr-1" md={2}>
        <Form.Group>
          <label>Recurring</label>
            <Form.Check
              label={"Daily"}
              value={"Hour"}
              type={"radio"}
              id={`1-recurring-${props.id}`}
              checked={props.rule.recurring === "Hour"}
              name={"recurring"}
              onChange={e => props.handleRuleChange(e, props.id)}
            />
            <Form.Check
              label={"Hourly"}
              value={"Minute"}
              type={"radio"}
              id={`2-recurring-${props.id}`}
              checked={props.rule.recurring === "Minute"}
              name={"recurring"}
              onChange={e => props.handleRuleChange(e, props.id)}
            />
            <Form.Check
              label={"By Minute"}
              value={"Second"}
              type={"radio"}
              id={`3-recurring-${props.id}`}
              checked={props.rule.recurring === "Second"}
              name={"recurring"}
              onChange={e => props.handleRuleChange(e, props.id)}
            />
        </Form.Group>
      </Col>
      <Col className="pr-1" md="3">
        <Form.Group>
          <label>START ON THE - {props.rule.recurring}</label>
          <Form.Control
            id={"start"}
            name={"start"}
            readOnly={props.rule.type === "image"}
            value={props.rule.start}
            placeholder="start"
            type="text"
            onChange={e => props.handleRuleChange(e, props.id)}
          >
          </Form.Control>
        </Form.Group>
      </Col>
      <Col className="pl-1" md="3">
        <Form.Group>
          <label>Duration (sec)/ #pictures</label>
          <Form.Control
            id={"duration"}
            name={"duration"}
            value={props.rule.duration}
            placeholder="duration"
            type="text"
            onChange={e => props.handleRuleChange(e, props.id)} />
        </Form.Group>
      </Col>
      <Col className="text-center" md="1">
        <Form.Group>
          <Button
            size={"sm"}
            className={"btn-fill pull-right"}
            onClick={e => props.removeRule(e, props.id)}
          > - </Button>
        </Form.Group>
      </Col>
    </Row>
  )
}

export default RuleRow