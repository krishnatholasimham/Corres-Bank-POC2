#!/usr/bin/env ruby

require "json"
require "table_print"

def join_arrays row
  row.keys.each do |key|
    row[key] = row[key].join(', ') if row[key].respond_to? :join
  end
end

d = JSON.parse(ARGF.read).each { |row| join_arrays row }
w = `/usr/bin/env tput cols`.to_i
tp.set :max_width, w
tp d, "MsgNum","Type","OrderingInstitution","AccountWithInstitution","PaymentAmount","IndicativeBalance","StatementAccountBalance","UnconfirmedSortTime","Sequence","LocalTime"
