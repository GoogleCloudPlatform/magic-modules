fun replaceCharsId(id: String): String{
    var newId = id.replace("-", "")
    newId = newId.replace(" ", "_")
    newId = newId.uppercase()

    return newId
}