require "gitm_pb"

module Gitm
  module Protobuf
    class Object
      def inspect
        "<#{self.class}: hash: #{self["hash"].unpack("H*").first}>"
      end
    end
  end
end
