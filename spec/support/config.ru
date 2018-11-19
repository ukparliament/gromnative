app = Proc.new do |env|
  not_found_response = ['404', {'Content-Type' => 'text/html'}, ['File not found']]
  found_response     = ['200', {'Content-Type' => 'application/ntriples'}]
  error_response     = ['500', {'Content-Type' => 'text/html'}, ['Server error']]
  redirect_response  = ['302', {'Content-Type' => 'text/html', 'Location' => 'http://localhost:3333/full.nt'}, ['Found elsewhere']]

  path = env['REQUEST_PATH']
  if path == '/500'
    error_response
  elsif path == '/404'
    not_found_response
  elsif path == '/302'
    redirect_response
  else
    fixture_path = "../fixtures#{path}"

    file_exists = File.file?(fixture_path)
    if file_exists
      fixture_data = File.read(fixture_path)
      found_response << [fixture_data]
    else
      not_found_response
    end
  end
end

run app