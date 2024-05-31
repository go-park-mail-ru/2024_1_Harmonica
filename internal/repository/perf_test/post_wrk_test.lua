local wrk = require("wrk")

-- Функция для чтения файла
local function read_file(path)
  local file, errorMessage = io.open(path, "rb")
  if not file then
    error(errorMessage)
  end
  local content = file:read("*all")
  file:close()
  return content
end

local boundary = "----WebKitFormBoundaryePkpFF7tjBAqx29L"

wrk.method = "POST"
wrk.path = "/api/v1/pins"
wrk.headers["Content-Type"] = "multipart/form-data; boundary=" .. boundary

local fileBody = read_file("image.jpeg")
local contentDispositionImage = 'Content-Disposition: form-data; name="image"; filename="image.jpeg"'
local contentTypeImage = 'Content-Type: image/jpeg'
local filePart = string.format(
  '--%s\r\n%s\r\n%s\r\n\r\n%s\r\n',
  boundary, contentDispositionImage, contentTypeImage, fileBody
)

local title = tostring(math.random(0, 100000))
local description = tostring(math.random(0, 100000))
local pinData = string.format('{"title":"%s", "description":"%s"}', title, description)
local contentDispositionPin = 'Content-Disposition: form-data; name="pin"'
local contentTypePin = 'Content-Type: application/json'
local pinPart = string.format(
  '--%s\r\n%s\r\n%s\r\n\r\n%s\r\n',
  boundary, contentDispositionPin, contentTypePin, pinData
)

local body = filePart .. pinPart .. '--' .. boundary .. '--\r\n'
wrk.body = body
