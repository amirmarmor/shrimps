import React, {useEffect, useState} from "react"
import {useDispatch, useSelector} from "react-redux"
import {selectConfig, setConfigAsync} from "../features/config/configSlice"
import {Button, Card, Col, Container, Form, Row,} from "react-bootstrap"
import RuleRow from "../components/Configuration/rules"
import Instructions from "../components/Configuration/instructions"

function Configuration() {
  const config = useSelector(selectConfig)
  const dispatch = useDispatch()
  const [currentConfig, setConfig] = useState()

  useEffect(() => {
    setConfig(config)
  }, [config])

  function addRule() {
    let id = 0
    if (!currentConfig) {
      return
    }
    console.log(currentConfig)
    if (currentConfig && currentConfig.rules.length > 0) {
      id = parseInt(currentConfig.rules[currentConfig.rules.length - 1].id) + 1
    }
    let rule = {
      id,
      recurring: 'Hour',
      start: 0,
      duration: 0
    }
    setConfig({
      ...currentConfig,
      rules: [
        ...currentConfig.rules,
        rule
      ]
    })
  }

  function removeRule(e, ind) {
    let rules = [...currentConfig.rules]
    rules.splice(ind, 1)
    setConfig({
      ...currentConfig,
      rules,
    })
  }

  function handleRuleChange(e, id) {
    console.log(e.target, id, currentConfig.rules)
    let rules = [...currentConfig.rules]
    rules[id] = {
      ...rules[id],
      [e.target.name]: e.target.value
    }
    setConfig({
      ...currentConfig,
      rules
    })
  }

  function handleChange(e) {
    let value = e.target.value
    if(e.target.name === "cleanup"){
      value = !currentConfig.cleanup
    }

    setConfig({
      ...currentConfig,
      [e.target.name]: value
    })
  }

  function handleSubmit(e) {
    e.preventDefault()
    let rules = currentConfig.rules.map(rule => {
      return {
        ...rule,
        id: rule.id.toString(),
        start: rule.start.toString(),
        duration: rule.duration.toString(),
      }
    })
    console.log("submit", currentConfig, rules)
    dispatch(setConfigAsync({rules, offset: currentConfig.offset.toString(), cleanup: currentConfig.cleanup}))
  }

  function renderRules() {
    if (currentConfig && currentConfig.rules) {
      return currentConfig.rules.map((rule, i) =>
        <RuleRow
          rule={rule}
          handleRuleChange={handleRuleChange}
          removeRule={removeRule}
          id={i}
          key={`rule-row-${i}`}
        />
      )
    }
    return ''
  }

  function getCleanupValue(){
    if(!currentConfig){
      return false
    } else {
      return currentConfig.cleanup
    }
  }

  return (
    <Container fluid>
      <Row>
        <Col md="8">
          <Card>
            <Card.Header>
              <Card.Title as="h4">Edit Configuration</Card.Title>
            </Card.Header>
            <Card.Body>
              <Form onSubmit={handleSubmit}>
                <Row>
                  <Col className="px-1" md="5">
                    <Form.Group>
                      <label>Offset</label>
                      <Form.Control
                        value={currentConfig ? currentConfig.offset : '0'}
                        placeholder="Offset"
                        type="text"
                        name={"offset"}
                        onChange={e => handleChange(e)}
                      />
                    </Form.Group>
                  </Col>
                  <Col className="px-1" md="5">
                    <label>Clean up</label>
                    < br/>
                    <input
                      name={"cleanup"}
                      id={"cleanup"}
                      type={"checkbox"}
                      value={getCleanupValue()}
                      checked={getCleanupValue()}
                      onChange={e => handleChange(e)}
                    />
                  </Col>
                </Row>
                <Row>
                  <Card.Header>Rules</Card.Header>
                </Row>
                {renderRules()}
                <Button
                  className={"btn-fill pull-right"}
                  style={{marginRight: "10px"}}
                  onClick={addRule}
                > + </Button>
                <Button
                  className="btn-fill pull-right"
                  type="submit"
                  variant="info"
                >
                  Update Configuration
                </Button>
                <div className="clearfix"/>
              </Form>
            </Card.Body>
          </Card>
        </Col>
        <Col md="4">
          <Instructions/>
        </Col>
      </Row>
    </Container>
  )
}

export default Configuration
