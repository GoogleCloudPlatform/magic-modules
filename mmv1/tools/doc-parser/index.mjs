import {JSDOM} from "jsdom"
import fetch from "node-fetch"
import fs from "fs"
import TurndownService from 'turndown'
var turndownService = new TurndownService()


const host = "https://cloud.google.com"
const path = "/dlp/docs/reference/rest/v2/projects.deidentifyTemplates#DeidentifyTemplate.CryptoReplaceFfxFpeConfig"
const subObject = path.split('#').pop().split(".").join("\\.")
const responses = {}

async function fetchUrlContents(path){
  var url = host + path
  var absPath = path.split("#")[0]
  if( !responses[absPath]) {
    responses[absPath] = await (await fetch(url)).text()
  }
  return responses[absPath]
}

function htmlToMarkDown(element) {
  return turndownService.turndown(element.outerHTML)
}

function analyzeDescription(elements) {
  elements = [... elements]
  var description = elements.map( e => htmlToMarkDown(e)).join('\n\n').trim()
  var required = false
  if (description.indexOf("Required.") == 0 ){
    description = description.substring('Required.'.length)
    required = true
  }
  return {description, required}
}

function printLine(tabIndex, line){
  return "  ".repeat(tabIndex) + line + '\n'
}

function printDescription(tabIndex, description){
  return printLine(
    tabIndex, "description: | " + "\n" +
    "  ".repeat(tabIndex+1) + description.split("\n").map(s=>s.trim()).join("\n" + "  ".repeat(tabIndex+1)).trim()
  )
}

function shouldBeFiltered(key) {
  var oneof = ['oneof_source', 'oneof_alphabet', 'oneof_characters', "oneof_transformation", "oneof_type"]
  return oneof.findIndex(v => key.indexOf(v) != -1) == -1
}

async function generateEnumValues(path, object, tabIndex) {
  var text = await fetchUrlContents(path)
  var {window} = await new JSDOM(text);
  var fields
  if (!object) {
    fields = window.document.querySelectorAll(`tr[id^=ENUM_VALUES]`)
  } else {
    fields = window.document.querySelectorAll(`tr[id^=${object}\\.ENUM_VALUES]`)
  }
  fields = [... fields]
  fields.filter(field => field.getAttribute('id').indexOf('UNSPECIFIED') == -1)
  var out = await Promise.all(fields.map( async field => {
    var row = field.querySelectorAll('td')
    row = [... row]
    return "  ".repeat(tabIndex+1) + `- :${row[0].textContent.trim()}  #${row[1].textContent.trim().replaceAll("\n","")}`
  }))
  return out.join("\n")
}

async function generateProperties(window, object, tabIndex){
  var fields
  if (!object) {
    fields = window.document.querySelectorAll(`tr[id^=FIELDS]`)
  } else {
    fields = window.document.querySelectorAll(`tr[id^=${object}\\.FIELDS]`)
  }
  fields = [... fields]
  fields = fields.filter( element => shouldBeFiltered(element.getAttribute("id")))
  var out = await Promise.all(fields.map( async element => {
    console.log(element.getAttribute("id"))
    var row = [... element.querySelectorAll(`td`)]
    var name = row[0].textContent
    if (row[1] == undefined) {
      console.log(name)
      console.log(object)
    }
    var val = row[1].children
    val = [... val]
    var inter = val.shift()
    var descriptionAnalyzed = analyzeDescription(val)
    var description = descriptionAnalyzed.description
    var required = descriptionAnalyzed.required
    if(inter.textContent.indexOf("object") != -1){
      var newPath = inter.querySelector('a').getAttribute('href')
      var newObject = newPath.split('#').length > 1 ?
       newPath.split('#').pop().split(".").join("\\."):
       ""
      return await buildObject(newPath, newObject, tabIndex, name, description, required)
    }
    var sout
    if(inter.textContent.indexOf("string") != -1){
      sout = printLine(tabIndex, "- !ruby/object:Api::Type::String")
    } else if(inter.textContent.indexOf("number") != -1){
      sout = printLine(tabIndex, "- !ruby/object:Api::Type::Double")
    } else if(inter.textContent.indexOf("boolean") != -1){
      sout = printLine(tabIndex, "- !ruby/object:Api::Type::Boolean")
    } else if(inter.textContent.indexOf("integer") != -1){
      sout = printLine(tabIndex, "- !ruby/object:Api::Type::Integer")
    } else if(inter.textContent.indexOf("enum") != -1){
      sout = printLine(tabIndex, "- !ruby/object:Api::Type::Enum")
    } else {
      sout = printLine(tabIndex, "- !WARNING_NOT_SUPPORTED " + inter.textContent)
    }
    var forkedTabIndex = tabIndex + 1
    sout += printLine(forkedTabIndex, `name: '${name}'`)
    if (required){
      sout += printLine(forkedTabIndex,"required: true")
    }
    sout += printDescription(forkedTabIndex, description)
    if(inter.textContent.indexOf("enum") != -1){
      sout += printLine(forkedTabIndex, `values:`)
      var newPath = inter.querySelector('a').getAttribute('href')
      var newObject = newPath.split('#').length > 1 ?
        newPath.split('#').pop().split(".").join("\\."):
        ""
      sout += await generateEnumValues(newPath, newObject, forkedTabIndex)
      sout += "\n"
    }
    forkedTabIndex--
    return sout
  }))
  return out.join('')
}

async function buildObject(path, object ,tabIndex, name, description, required){
  var text = await fetchUrlContents(path)
  var {window} = await new JSDOM(text);
  var name = name || object.split(".").pop()

  if(!description){
    var analyzedDescript = analyzeDescription(window.document.querySelector(`#${object}\\.description`).children)
    description = analyzedDescript.description
    required = analyzedDescript.required
  }
  var out = printLine(tabIndex, "- !ruby/object:Api::Type::NestedObject")
  tabIndex++
  out += printLine(tabIndex, `name: '${name}'`)
  out += printDescription(tabIndex, description)
  out += printLine(tabIndex, `properties:`)
  tabIndex++
  out += await generateProperties(window, object, tabIndex)
  return out
}



var out = await buildObject(path, subObject, 0)
fs.writeFileSync("./out.yaml", out)