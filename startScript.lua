#!/usr/bin/env tarantool


box.cfg{
}


box.schema.user.grant('guest', 'read,write,execute', 'universe')
box.schema.user.passwd('admin', 'admin')


---------------------------------------



s = box.schema.space.create('tester', { if_not_exists = true, engine = 'memtx' })


s:format({
                { name = 'id', type = 'unsigned' },
                { name = 'advertiser_id', type = 'unsigned' },
                { name = 'static_link', type = 'string' },
                { name = 'is_active', type = 'boolean' },
                { name = 'limit', type = 'unsigned' },
                { name = 'data_start', type = 'unsigned' },
                { name = 'data_end', type = 'unsigned' },
                { name = 'target', type = 'unsigned' },

            })


s:create_index('primary', { type = 'tree', parts = { 'id' } })
s:create_index('secondary', { unique = true, type = 'hash', parts = { 'target' } })


-----------------------------------

t = box.schema.space.create('target', { if_not_exists = true, engine = 'memtx' })

t:format({
                { name = 'id', type = 'unsigned' },
                { name = 'advertiser_id', type = 'unsigned' },
                { name = 'theme', type = 'unsigned' },
                { name = 'gender', type = 'unsigned' },
                { name = 'age', type = 'unsigned' },
                { name = 'geographic_id', type = 'unsigned' },
                { name = 'wealth', type = 'unsigned' },
            })

t:create_index('primary_target', { unique = true, type = 'tree', parts = { 'id' } })


---------------------------



function myget(id_adver, idIndex)
    val = box.space.tester:select{idIndex}
    return box.execute( 'select * from "tester" where "advertiser_id" in (' .. id_adver .. ');' )
end



-- function myget(id_adver, idIndex)
--     val = box.space.tester:select{idIndex}
--     return box.execute( 'select * from "target" where "theme" in (' .. id_adver .. ')
--     AND "gender" = true AND "age" > 0 AND "geographic_id"  in (' .. idIndex .. ');' )
-- end