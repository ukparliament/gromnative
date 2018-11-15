app = Proc.new do |env|
  not_found_response = ['404', {'Content-Type' => 'text/html'}, ['File not found']]
  found_response = ['200', {'Content-Type' => 'application/ntriples'}]

  path = env['REQUEST_PATH']
  fixture_path = "../fixtures#{path}"

  file_exists = File.file?(fixture_path)
  if file_exists
    fixture_data = File.read(fixture_path)
    found_response << [fixture_data]
  else
    not_found_response
  end
end

run app