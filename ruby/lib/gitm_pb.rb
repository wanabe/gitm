# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: gitm.proto

require 'google/protobuf'

Google::Protobuf::DescriptorPool.generated_pool.build do
  add_file("gitm.proto", :syntax => :proto3) do
    add_message "gitm.protobuf.Object" do
      optional :hash, :string, 1
    end
    add_message "gitm.protobuf.Repository" do
      optional :path, :string, 1
    end
    add_message "gitm.protobuf.Commit" do
      optional :object, :message, 1, "gitm.protobuf.Object"
      repeated :parents, :message, 2, "gitm.protobuf.Object"
    end
    add_message "gitm.protobuf.LogRequest" do
      optional :repository, :message, 1, "gitm.protobuf.Repository"
      optional :object, :message, 2, "gitm.protobuf.Object"
    end
  end
end

module Gitm
  module Protobuf
    Object = Google::Protobuf::DescriptorPool.generated_pool.lookup("gitm.protobuf.Object").msgclass
    Repository = Google::Protobuf::DescriptorPool.generated_pool.lookup("gitm.protobuf.Repository").msgclass
    Commit = Google::Protobuf::DescriptorPool.generated_pool.lookup("gitm.protobuf.Commit").msgclass
    LogRequest = Google::Protobuf::DescriptorPool.generated_pool.lookup("gitm.protobuf.LogRequest").msgclass
  end
end
