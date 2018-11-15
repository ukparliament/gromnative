RSpec.describe GromNative do
  it 'has a version number' do
    expect(GromNative::VERSION).not_to be nil
  end

  describe '.get' do
    context 'with a valid url' do
      it 'makes the expected request' do
        val = subject.get('http://localhost:3333/full.nt')

        puts val
      end
    end

    context 'with an invalid url' do
      it 'returns the expected object' do
        expect(JSON.parse(subject.get('foo://a_broken.url'))).to eq({'statementsBySubject' => nil,'edgesBySubject' => nil,'statusCode' => 0,'uri' => '','error' => "Error getting data: Get foo://a_broken.url: unsupported protocol scheme \"foo\"\n"})
      end
    end
  end
end
